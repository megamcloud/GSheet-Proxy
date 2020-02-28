package memDb_test

import (
	"git.anphabe.net/event/anphabe-event-hub/domain/service"
	. "git.anphabe.net/event/anphabe-event-hub/infrastructure/repository/memDb"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("MemDbConnection", func() {

	Context("Call to NewMemDbConnection() to create a *memDb.Connection", func() {
		connection := NewMemDbConnection()

		It("connection must be *memDb.Connection", func() {
			Expect(connection).To(BeAssignableToTypeOf((*Connection)(nil)))
		})

		It("connection should implements service.DbConnectionInterface", func() {
			Expect(connection).To(BeAssignableToTypeOf((service.DbConnectionInterface)(connection)))
		})
	})

	Context("given *memDb.Connection created", func() {
		sut := NewMemDbConnection()
		repoName := "testRepo"

		Context("calling to InitRepository("+`"`+repoName+`"`+") to create a *ScanItemRepository", func() {

			repository, _ := sut.InitRepository(repoName)

			It("repository must be *memDb.ScanItemRepository", func() {
				Expect(repository).To(BeAssignableToTypeOf((*ScanItemRepository)(nil)))
			})

			It("repository's name must be "+`"`+repoName+`"`, func() {
				Expect(repository.GetRepoName()).To(Equal(repoName))
			})
		})
	})
})
