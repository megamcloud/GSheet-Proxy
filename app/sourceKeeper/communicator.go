package sourceKeeper

import (
	"errors"
	"fmt"
	"git.anphabe.net/event/anphabe-event-hub/config"
	"git.anphabe.net/event/anphabe-event-hub/domain/model/dbSource"
	"git.anphabe.net/event/anphabe-event-hub/infrastructure/jsonapiClient"
	"go.uber.org/zap"
)

type Communicator struct {
	name          string
	idField       string
	source        *dbSource.DbSource
	logger        *zap.Logger
	importRunning bool
}

func NewCommunicator(cfg config.DbSource, logger *zap.Logger) *Communicator {
	return &Communicator{
		name:          cfg.Name,
		idField:       cfg.IdField,
		source:        dbSource.NewDBSource(cfg),
		importRunning: false,
		logger:        logger,
	}
}

func (i *Communicator) Import(callback func(repoName string, idField string, data []map[string]string) int) error {
	if i.importRunning {
		return errors.New("there is an import currently importRunning")
	} else {
		i.importRunning = true
	}

	defer func() { i.importRunning = false }()
	i.logger.Info("Communicator: start import", zap.String("dbName", i.name))

	progress := newFetchingProgress(i.source.GetFetchingUrl(0, 200))

	var count int = 0

	for progress.HasNext() {
		if data, err := progress.FetchNext(); err == nil {
			count += callback(i.name, i.idField, data)
		} else {
			i.logger.Error("Communicator error: "+err.Error(), zap.String("dbName", i.name))
			return err
		}
	}

	i.logger.Info("Communicator: finish import", zap.String("dbName", i.name), zap.Int("numItem", count))

	return nil
}

func (i *Communicator) Update(key string, params map[string]string) bool {
	updateUrl := i.source.GetUpdateUrl(key, params)
	fmt.Println(updateUrl)

	if _, err := jsonapiClient.Get(updateUrl); nil == err {
		fmt.Print(updateUrl)
		return true
	} else {
		fmt.Println(err.Error())
	}

	return false
}
