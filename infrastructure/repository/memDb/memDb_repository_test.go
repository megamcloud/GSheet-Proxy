package memDb_test

import (
	"errors"
	"fmt"
	"git.anphabe.net/event/anphabe-event-hub/domain/model/scanItem"
	. "git.anphabe.net/event/anphabe-event-hub/infrastructure/repository/memDb"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"time"
)

var _ = Describe("*memDb.ScanItemRepository", func() {

	Describe(":: GetRepoName() to get Repository's name", func() {
		repoName := "testRepo"
		repo, _ := NewMemDbConnection().InitRepository(repoName)
		defer repo.CloseDb()

		name := repo.GetRepoName()

		It("return must be "+`"`+repoName+`"`, func() {
			Expect(name).To(Equal(repoName))
		})
	})

	////////////////////////////////////////////////////////////////////////////

	Describe(":: GetItem(key) to get an Item", func() {
		Context(" :GIVEN: an item was-not-exist", func() {

			repoName := "testRepo"
			repo, _ := NewMemDbConnection().InitRepository(repoName)
			defer repo.CloseDb()

			Context(" :THEN ⇶ call to GetItem(key-not-exist)", func() {

				got, found := repo.GetItem("key-does-not-exist")

				It("found should be FALSE", func() {
					Expect(found).To(BeFalse())
				})

				It("got should be NIL", func() {
					Expect(got).To(BeNil())
				})
			})
		})

		Context(" :GIVEN: an item was-exist", func() {
			repoName := "testRepo"
			repo, _ := NewMemDbConnection().InitRepository(repoName)
			defer repo.CloseDb()

			item := scanItem.CreateTestScanItem("1a")
			repo.SetItem(item)

			Context(" :THEN ⇶ call to GetItem(key-exist)", func() {
				got, found := repo.GetItem(item.GetKey())

				It("found should be FALSE", func() {
					Expect(found).To(BeTrue())
				})

				It("got should be NIL", func() {
					Expect(got).To(Equal(item))
				})
			})
		})
	})

	////////////////////////////////////////////////////////////////////////////

	Describe(":: SetItem(item) to save an Item", func() {
		Context(" :GIVEN: an item was-not-exist", func() {

			repoName := "testRepo"
			repo, _ := NewMemDbConnection().InitRepository(repoName)
			defer repo.CloseDb()

			Context(" :THEN ⇶ call to SetItem(item)", func() {
				item := scanItem.CreateTestScanItem("1a")
				repo.SetItem(item)
				got, found := repo.GetItem(item.GetKey())

				It("item must be successful kept in repository", func() {
					Expect(got).To(Equal(item))
					Expect(found).To(BeTrue())
				})
			})
		})

		Context(" :GIVEN: an item was-exist", func() {
			repoName := "testRepo"
			repo, _ := NewMemDbConnection().InitRepository(repoName)
			defer repo.CloseDb()

			item1 := scanItem.CreateTestScanItem("1a")
			repo.SetItem(item1)
			got1, _ := repo.GetItem(item1.GetKey())

			Context(" :THEN ⇶ call to SetItem(item) to new item with and old-key-existed", func() {
				item2 := scanItem.CreateTestScanItem("1b")
				repo.SetItem(item2)

				got2, found := repo.GetItem(item2.GetKey())

				It("new-item must overwrite old-item", func() {
					Expect(found).To(BeTrue())
					Expect(got2).To(Equal(item2))
					Expect(got2).NotTo(Equal(got1))
				})
			})
		})
	})

	////////////////////////////////////////////////////////////////////////////

	Describe(":: NewItem(Key, Data) to save an Item", func() {
		repoName := "testRepo"
		sut, _ := NewMemDbConnection().InitRepository(repoName)

		AfterEach(func() {
			sut.CloseDb()
		})

		Context(" :GIVEN: an item-key was EMPTY/BLANK", func() {
			var inputKey string = ""
			var inputData = map[string]string{"field1": "data1", "field2": "data2"}

			Context(" :THEN ⇶ call to NewItem(inputKey, inputData)", func() {
				got, err := sut.NewItem(inputKey, inputData)

				It("error will be returned", func() {
					Expect(err).To(Equal(errors.New("key could not be empty")))
					Expect(got).To(BeNil())
				})
			})
		})

		Context(" :GIVEN: an item-key valid", func() {
			repoName := "testRepo"
			sut, _ := NewMemDbConnection().InitRepository(repoName)

			AfterEach(func() {
				sut.CloseDb()
			})

			Context(" :THEN ⇶ call to SetItem(item) to new item with and old-key-existed", func() {
				item := scanItem.CreateTestScanItem("1a")
				got1, err := sut.NewItem(item.GetKey(), item.GetData())
				got2, found := sut.GetItem(item.GetKey())

				It("item success save in repository", func() {
					Expect(err).To(BeNil())
					Expect(found).To(BeTrue())
					Expect(got1).To(Equal(item))
					Expect(got1).To(Equal(got2))
				})
			})
		})
	})

	////////////////////////////////////////////////////////////////////////////

	Describe(":: AddActivities(itemKey, action, properties) to add new Item's Activities", func() {
		var repoName string
		var itemKey string
		var repo scanItem.RepositoryInterface
		var fakeActivity1, fakeActivity2 scanItem.ItemActivity

		BeforeEach(func() {
			repoName = "testRepo"
			itemKey = "testItem"
			repo, _ = NewMemDbConnection().InitRepository(repoName)
			_, _ = repo.NewItem(itemKey, map[string]string{"field1": "data1", "field2": "data2"})

			fakeActivity1 = scanItem.ItemActivity{
				Action:  "testAction1",
				Data:    map[string]string{"field1": "data1", "field2": "data2"},
				Created: time.Now(),
			}

			fakeActivity2 = scanItem.ItemActivity{
				Action:  "testAction2",
				Data:    map[string]string{"field1": "data1", "field2": "data2"},
				Created: time.Now(),
			}
		})

		AfterEach(func() {
			repo.CloseDb()
		})

		Context(" :GIVEN: an item with no Activity", func() {
			Context(" :THEN ⇶ call to AddActivities(itemKey)", func() {
				It("should return correct *ItemActivities", func() {
					got := repo.AddItemActivity(itemKey, fakeActivity1)

					By("by return type *ItemActivities")
					Expect(got).To(BeAssignableToTypeOf((*scanItem.ItemActivities)(nil)))

					By("by return correct itemKey")
					Expect(got.Key).To(Equal(itemKey))

					By("by return 01 activity")
					Expect(len(got.Activities)).To(Equal(1))

					By("by return the same activity data")
					Expect(got.Activities[0]).To(Equal(fakeActivity1))
				})
			})
		})

		Context(" :GIVEN: an item with 1 Activity", func() {
			Context(" :THEN ⇶ call to AddActivities(itemKey)", func() {
				It("should return correct *ItemActivities", func() {
					_ = repo.AddItemActivity(itemKey, fakeActivity1)
					got := repo.AddItemActivity(itemKey, fakeActivity2)

					By("by return type *ItemActivities")
					Expect(got).To(BeAssignableToTypeOf((*scanItem.ItemActivities)(nil)))

					By("by return correct itemKey")
					Expect(got.Key).To(Equal(itemKey))

					By("by return 01 activity")
					Expect(len(got.Activities)).To(Equal(2))

					By("latest added activity should be on top")
					fmt.Printf("%v", got.Activities)
					Expect(got.Activities[0]).To(Equal(fakeActivity2))
					Expect(got.Activities[1]).To(Equal(fakeActivity1))

				})
			})
		})
	})

	////////////////////////////////////////////////////////////////////////////

	Describe(":: GetActivities(Key, Data) to get an Item's Activities", func() {
		var repoName string
		var itemKey string
		var repo scanItem.RepositoryInterface
		var fakeActivity1, fakeActivity2 scanItem.ItemActivity

		BeforeEach(func() {
			repoName = "testRepo"
			itemKey = "testItem"
			repo, _ = NewMemDbConnection().InitRepository(repoName)
			_, _ = repo.NewItem(itemKey, map[string]string{"field1": "data1", "field2": "data2"})

			fakeActivity1 = scanItem.ItemActivity{
				Action:  "testAction1",
				Data:    map[string]string{"field1": "data1", "field2": "data2"},
				Created: time.Now(),
			}

			fakeActivity2 = scanItem.ItemActivity{
				Action:  "testAction2",
				Data:    map[string]string{"field1": "data1", "field2": "data2"},
				Created: time.Now(),
			}
		})

		AfterEach(func() {
			repo.CloseDb()
		})

		Context(" :GIVEN: no item exist", func() {
			Context(" :THEN ⇶ call to GetActivities(itemKey)", func() {
				It("should return nil", func() {
					got := repo.GetItemActivities("not exist")
					Expect(got).To(BeNil())
				})
			})
		})

		Context(" :GIVEN: an item with no Activity", func() {
			Context(" :THEN ⇶ call to GetActivities(itemKey)", func() {
				It("should return correct *ItemActivities", func() {
					got := repo.GetItemActivities(itemKey)

					By("by return type *ItemActivities")
					Expect(got).To(BeAssignableToTypeOf((*scanItem.ItemActivities)(nil)))

					By("by return correct itemKey")
					Expect(got.Key).To(Equal(itemKey))

					By("should return 0 activity", func() {
						Expect(got.Activities).To(BeNil())
						Expect(len(got.Activities)).To(Equal(0))
					})
				})
			})
		})

		Context(" :GIVEN: an item with 01 Activity", func() {
			Context(" :THEN ⇶ call to GetActivities(itemKey)", func() {
				It("should return correct *ItemActivities", func() {
					_ = repo.AddItemActivity(itemKey, fakeActivity1)
					got := repo.GetItemActivities(itemKey)

					By("by return type *ItemActivities")
					Expect(got).To(BeAssignableToTypeOf((*scanItem.ItemActivities)(nil)))

					By("by return correct itemKey")
					Expect(got.Key).To(Equal(itemKey))

					By("by return 01 activity")
					Expect(len(got.Activities)).To(Equal(1))

					By("by return the same activity data")
					Expect(got.Activities[0]).To(Equal(fakeActivity1))
				})
			})
		})

		Context(" :GIVEN: an item with 02 Activities", func() {
			Context(" :THEN ⇶ call to GetActivities(itemKey)", func() {
				It("should return correct *ItemActivities", func() {
					_ = repo.AddItemActivity(itemKey, fakeActivity1)
					_ = repo.AddItemActivity(itemKey, fakeActivity2)
					got := repo.GetItemActivities(itemKey)

					By("by return type *ItemActivities")
					Expect(got).To(BeAssignableToTypeOf((*scanItem.ItemActivities)(nil)))

					By("by return correct itemKey")
					Expect(got.Key).To(Equal(itemKey))

					By("by return 01 activity")
					Expect(len(got.Activities)).To(Equal(2))

					By("latest added activity should be on top")
					fmt.Printf("%v", got.Activities)
					Expect(got.Activities[0]).To(Equal(fakeActivity2))
					Expect(got.Activities[1]).To(Equal(fakeActivity1))

				})
			})
		})
	})

	////////////////////////////////////////////////////////////////////////////

	Describe(":: Len() to get number of items in repository", func() {
		Context(" :GIVEN: 2 items are put in repository", func() {
			repoName := "testRepo"
			sut, _ := NewMemDbConnection().InitRepository(repoName)
			item1 := scanItem.CreateTestScanItem("1a")
			item2 := scanItem.CreateTestScanItem("2")
			sut.SetItem(item1)
			sut.SetItem(item2)

			AfterEach(func() {
				sut.CloseDb()
			})

			Context(" :THEN ⇶ call to Len() ", func() {
				got := sut.Len()

				It(" should return 2", func() {
					Expect(got).To(Equal(2))
				})
			})
		})
	})

	////////////////////////////////////////////////////////////////////////////

	Describe(":: Items() to get all items in repository", func() {
		Context(" :GIVEN: 2 items are put in repository", func() {
			repoName := "testRepo"
			sut, _ := NewMemDbConnection().InitRepository(repoName)
			item1 := scanItem.CreateTestScanItem("1a")
			item2 := scanItem.CreateTestScanItem("2")
			sut.SetItem(item1)
			sut.SetItem(item2)

			AfterEach(func() {
				sut.CloseDb()
			})

			Context(" :THEN ⇶ call to Items() ", func() {
				got := sut.Items()

				It("will return all 2 items", func() {
					Expect(got).To(ContainElement(item1))
					Expect(got).To(ContainElement(item2))
				})
			})
		})
	})
})
