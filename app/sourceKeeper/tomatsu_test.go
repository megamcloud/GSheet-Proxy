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
//)
//
//var _ = Describe("Tomatsu Real Test\n", func() {
//
//	fakeDbSource := config.DbSource{
//		Name:           "Tomatsu",
//		IdField:        "field_qrcode",
//		FetchingUrl:    "https://script.google.com/macros/s/AKfycbyKxlzZMiVlF01ZGPAXYsY0ARV-L8V04QCgONo5kIbTAwkfOC4C/exec?path=/sample&order=field_qrcode&offset=%offset%&limit=%size%",
//		FetchingFormat: "json",
//		UpdateUrl:      "",
//		UpdateMethod:   "",
//	}
//
//	Context(fmt.Sprintf("Given a real Tomatsu dbSource at %s", fakeDbSource.FetchingUrl), func() {
//
//		memRegistry := service.NewRepositoryRegistry(memDb.NewMemDbConnection())
//		repo, _ := memRegistry.GetRepository(fakeDbSource.Name)
//
//		Context("Calling to Import()\n", func() {
//			fakeLogger, _ := fakeLogger()
//			_ = sourceKeeper.NewCommunicator(fakeDbSource, fakeLogger)
//			//_ = importer.Import()
//			//
//			//FIt("Every items must be saved in repository\n", func() {
//			//	Expect(repo.Len()).To(Equal(611))
//			//})
//		})
//	})
//
//})
