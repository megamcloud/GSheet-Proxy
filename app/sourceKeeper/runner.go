package sourceKeeper

import (
	"fmt"
	"git.anphabe.net/event/anphabe-event-hub/config"
	"git.anphabe.net/event/anphabe-event-hub/domain/model/scanItem"
	"git.anphabe.net/event/anphabe-event-hub/domain/service"
	"go.uber.org/zap"
	"strings"
	"sync"
	"time"
)

const SyncTime time.Duration = 5 * time.Minute

type Keeper struct {
	conf         []config.DbSource
	repoRegistry service.RepositoryRegistryInterface
	dbSources    map[string]*Communicator
	nextImports  map[string]time.Time
	stop         chan struct{}
	importChan   chan string
	activityChan chan *activityLog
	logger       *zap.Logger
}

type activityLog struct {
	repoName string
	itemKey  string
	activity scanItem.ItemActivity
}

func NewSourceKeeper(cfg []config.DbSource, registry service.RepositoryRegistryInterface, logger *zap.Logger) *Keeper {

	instance := &Keeper{
		conf: cfg,
		repoRegistry: registry,
		dbSources:    make(map[string]*Communicator),
		nextImports:  make(map[string]time.Time),
		stop:         make(chan struct{}, 1),
		activityChan: make(chan *activityLog, 30),
		importChan:   make(chan string, 2),
		logger:       logger,
	}

	return instance
}

func (i *Keeper) Stop() {
	i.stop <- struct{}{}
}

func (i *Keeper) init() {
	for _, cfgDbSource := range i.conf {
		_, _ = i.repoRegistry.GetRepository(cfgDbSource.Name)
		i.updateNextImport(cfgDbSource.Name)
		i.dbSources[cfgDbSource.Name] = NewCommunicator(cfgDbSource, i.logger)
	}
}

func (i *Keeper) Start(wg *sync.WaitGroup) {
	i.init()

	// start broker
	go func() {
		tick := time.NewTicker(1 * time.Minute)
		var wgChild sync.WaitGroup

		defer func() {
			tick.Stop()
			wgChild.Wait()

			close(i.stop)
			close(i.activityChan)
			close(i.importChan)

			i.repoRegistry.Shutdown()
			wg.Done()

			i.logger.Info("SourceKeeper Stopped")
		}()

		i.logger.Info("SourceKeeper Started")

		for {
			select {
			case <-i.stop:
				i.logger.Info("SourceKeeper going to stop")
				return

			case activity, ok := <-i.activityChan:
				if ok {
					wgChild.Add(1)
					go i.pushActivityToSource(*activity, &wgChild)
				}

			case _, ok := <-tick.C:
				if ok {
					if repoName := i.pickImportSource(); repoName != "" {
						i.importChan <- repoName
					}
				}

			case repoName, ok := <-i.importChan:
				if ok {
					if repoName != "" {
						wgChild.Add(1)
						go i.importFromSource(repoName, &wgChild)
					}
				}
			}
		}
	}()
}

func (i *Keeper) StartImport(repoName string) bool {
	if _, err := i.repoRegistry.GetRepository(repoName); nil == err {
		fmt.Println(repoName)
		i.importChan <- repoName
		return true
	}

	return false
}

func (i *Keeper) GetItemDetail(repoName string, itemKey string) (*scanItem.ItemDetail, bool) {
	return i.getRepository(repoName).GetItemDetail(itemKey)
}

func (i *Keeper) GetItems(repoName string) []*scanItem.ScanItem {
	return i.getRepository(repoName).Items()
}

func (i *Keeper) ScanItem(repoName string, itemKey string, activityName string, properties map[string]string) (*scanItem.ItemDetail, bool) {
	i.addItemActivity(repoName, itemKey, activityName, properties)
	return i.GetItemDetail(repoName, itemKey)
}

func (i *Keeper) addItemActivity(repoName string, itemKey string, action string, properties map[string]string) *scanItem.ItemActivities {
	if _, found := i.GetItemDetail(repoName, itemKey); found {
		activity := scanItem.NewActivity(action, properties)
		activityLog := &activityLog{
			repoName: repoName,
			itemKey:  itemKey,
			activity: activity,
		}

		i.activityChan <- activityLog

		return i.getRepository(repoName).AddItemActivity(itemKey, activity)
	}

	return nil
}

func (i *Keeper) pushActivityToSource(log activityLog, wg *sync.WaitGroup) {
	defer wg.Done()

	// convert action-moment into string
	properties := log.activity.Data
	properties[log.activity.Action] = log.activity.Created.Format("2 Jan 2006 15:04:05")

	dbSource := i.dbSources[log.repoName]
	finish := dbSource.Update(log.itemKey, properties)

	if finish {
		i.logger.Info("Successful push activity", zap.String("dbName", log.repoName), zap.String("itemKey", log.itemKey))
	} else {
		// send back to queue to re-process
		sendBack := log
		i.activityChan <- &sendBack
		i.logger.Error("Fail pushing activity", zap.String("dbName", log.repoName), zap.String("itemKey", log.itemKey))
	}
}

func (i *Keeper) importFromSource(repoName string, wg *sync.WaitGroup) {
	defer wg.Done()

	_ = i.dbSources[repoName].Import(i.saveItems)
	i.updateNextImport(repoName)
}

func (i *Keeper) pickImportSource() string {
	for name, t := range i.nextImports {
		now := time.Now()
		if now.After(t) {
			return name
		}
	}

	return ""
}

func (i *Keeper) saveItems(repoName string, idField string, data []map[string]string) int {
	repo := i.getRepository(repoName)
	count := 0

	for _, item := range data {
		if key, found := item[idField]; found {

			if key = strings.TrimSpace(key); key != "" {
				_, err := repo.NewItem(key, item)

				if nil != err {
					i.logger.Error("Communicator: could not new item", zap.String("dbName", repoName))
				} else {
					count += 1
				}
			} else {
				i.logger.Error("Communicator: key empty", zap.String("dbName", repoName))
			}
		}
	}

	return count
}

func (i *Keeper) updateNextImport(repoName string) {
	if nextRun, exist := i.nextImports[repoName]; exist {
		i.nextImports[repoName] = nextRun.Add(SyncTime)
		i.logger.Info("Next importing schedule", zap.String("dbName", repoName), zap.Time("time", i.nextImports[repoName]))
	} else {
		i.nextImports[repoName] = time.Now()
	}
}

func (i *Keeper) getRepository(repoName string) scanItem.RepositoryInterface {
	if repo, err := i.repoRegistry.GetRepository(repoName); nil == err {
		return repo
	} else {
		panic(err)
	}
}
