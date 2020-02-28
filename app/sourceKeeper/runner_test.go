package sourceKeeper_test
//
//import (
//	"fmt"
//	"git.anphabe.net/event/anphabe-event-hub/app/sourceKeeper"
//	"git.anphabe.net/event/anphabe-event-hub/config"
//	"git.anphabe.net/event/anphabe-event-hub/domain/service"
//	"git.anphabe.net/event/anphabe-event-hub/infrastructure/repository/memDb"
//	. "github.com/onsi/ginkgo"
//	. "github.com/onsi/gomega"
//	"go.uber.org/zap"
//	"go.uber.org/zap/zapcore"
//	"go.uber.org/zap/zaptest/observer"
//	"gopkg.in/h2non/gock.v1"
//	"net/http"
//	"sync"
//	"time"
//)
//
//var _ = Describe("Keeper\n", func() {
//	dbSource1 := getFakeDbSource1()
//	dbSource2 := getFakeDbSource2()
//
//	Context(fmt.Sprintf("Given 2 DBSources will be setup at %s and %s \n", dbSource1.Domain, dbSource2.Domain), func() {
//		defer gock.Off()
//
//		for _, tt := range dbSource1.Steps {
//			gock.New(dbSource1.Domain).
//				Get(dbSource1.Path).
//				MatchParams(tt.wantParams).
//				Reply(http.StatusOK).
//				BodyString(tt.givenBody)
//		}
//
//		for _, tt := range dbSource2.Steps {
//			gock.New(dbSource2.Domain).
//				Get(dbSource2.Path).
//				MatchParams(tt.wantParams).
//				Reply(http.StatusOK).
//				BodyString(tt.givenBody)
//		}
//
//		fakeDbSource1 := config.DbSource{
//			Name:           dbSource1.Name,
//			IdField:        dbSource1.IdField,
//			FetchingUrl:    dbSource1.Domain + dbSource1.Path + dbSource1.QueryString,
//			FetchingFormat: "json",
//			UpdateUrl:      "",
//			UpdateMethod:   "",
//		}
//
//		fakeDbSource2 := config.DbSource{
//			Name:           dbSource1.Name,
//			IdField:        dbSource1.IdField,
//			FetchingUrl:    dbSource1.Domain + dbSource1.Path + dbSource1.QueryString,
//			FetchingFormat: "json",
//			UpdateUrl:      "",
//			UpdateMethod:   "",
//		}
//
//		stubLogger, observedLogs := fakeLogger()
//
//		memRegistry := service.NewRepositoryRegistry(memDb.NewMemDbConnection())
//		runner := sourceKeeper.NewSourceKeeper([]config.DbSource{fakeDbSource1, fakeDbSource2}, memRegistry, stubLogger)
//
//		repo1, _ := memRegistry.GetRepository(fakeDbSource1.Name)
//		repo2, _ := memRegistry.GetRepository(fakeDbSource2.Name)
//
//		Context("Calling to Start()\n", func() {
//			var wg sync.WaitGroup
//			wg.Add(1)
//			runner.Start(&wg)
//			var count int = 0
//			startLogs := observedLogs.FilterField(zap.String("state", "start"))
//			stopLogs := observedLogs.FilterField(zap.String("state", "stop"))
//
//			tick := time.NewTicker(sourceKeeper.SyncTime)
//
//		OuterLoop:
//			for {
//				select {
//				case <-tick.C:
//
//					if count >= 2 {
//						break OuterLoop
//					}
//
//					if startLogs.Len() >= 2 {
//						logs := startLogs.TakeAll()
//						It("It should has 2 log lines: DbImporter start import", func() {
//							for _, log := range logs {
//								Expect(log.Message).To(Equal("DbImporter start import"))
//								Expect([]string{fakeDbSource1.Name, fakeDbSource2.Name}).To(ContainElement(ContainSubstring(log.ContextMap()["dbName"].(string))))
//							}
//						})
//					}
//
//					if stopLogs.Len() >= 2 {
//						logs := stopLogs.TakeAll()
//						It("It should has 2 log lines: DbImporter finish import", func() {
//							for _, log := range logs {
//								Expect(log.Message).To(Equal("DbImporter finish import"))
//								Expect([]string{fakeDbSource1.Name, fakeDbSource2.Name}).To(ContainElement(ContainSubstring(log.ContextMap()["dbName"].(string))))
//							}
//						})
//
//						It("Every items must be saved in repository\n", func() {
//							for _, tt := range dbSource1.Steps {
//								for _, data := range tt.wantReturn {
//									item, found := repo1.GetItem(data[fakeDbSource1.IdField])
//									Expect(found).To(BeTrue())
//									Expect(data).To(Equal(item.GetData()))
//								}
//							}
//
//							for _, tt := range dbSource2.Steps {
//								for _, data := range tt.wantReturn {
//									item, found := repo2.GetItem(data[fakeDbSource2.IdField])
//									Expect(found).To(BeTrue())
//									Expect(data).To(Equal(item.GetData()))
//								}
//							}
//						})
//
//						count += 1
//					}
//				}
//			}
//
//			runner.Stop()
//		})
//	})
//})
//
//
