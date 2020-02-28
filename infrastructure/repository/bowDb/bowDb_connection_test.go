package bowDb_test

import (
	"git.anphabe.net/event/anphabe-event-hub/domain/service"

	"math/rand"
	"strconv"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "git.anphabe.net/event/anphabe-event-hub/infrastructure/repository/bowDb"
)

var _ = Describe("BowDbConnection", func() {

	Context("Call to NewBowDbConnection() to create a *bowDb.Connection", func() {
		connection := setupDb()

		It("connection must be *bowDb.Connection", func() {
			Expect(connection).To(BeAssignableToTypeOf((*Connection)(nil)))
		})

		It("connection should implements service.DbConnectionInterface", func() {
			Expect(connection).To(BeAssignableToTypeOf((service.DbConnectionInterface)(connection)))
		})
	})

	Context("given *bowDb.Connection created", func() {
		sut := setupDb()
		repoName := "testRepo"

		Context("calling to InitRepository("+`"`+repoName+`"`+") to create a *ScanItemRepository", func() {

			repository, _ := sut.InitRepository(repoName)

			It("repository must be *bowDb.ScanItemRepository", func() {
				Expect(repository).To(BeAssignableToTypeOf((*ScanItemRepository)(nil)))
			})

			It("repository's name must be "+`"`+repoName+`"`, func() {
				Expect(repository.GetRepoName()).To(Equal(repoName))
			})
		})
	})
})

func setupDb() *Connection {
	rand.Seed(time.Now().UnixNano())

	return NewBowDbConnection("./bow_data_" + strconv.Itoa(rand.Int()))
}