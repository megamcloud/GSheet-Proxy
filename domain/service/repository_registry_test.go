package service

import (
	"git.anphabe.net/event/anphabe-event-hub/domain/model/scanItem"
	"git.anphabe.net/event/anphabe-event-hub/infrastructure/repository/memDb"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRepositoryManager_GetDB__Given__Repo_NotExist__Expect__Return_RepositoryObject(t *testing.T) {
	tests := []struct {
		name          string
		givenRepoName string
		wantRepoName  string
		wantInterface *scanItem.RepositoryInterface
		wantErr       error
	}{
		{
			name:          "Get first Repository",
			givenRepoName: "repo1",
			wantRepoName:  "repo1",
			wantInterface: (*scanItem.RepositoryInterface)(nil),
			wantErr:       nil,
		},
		{
			name:          "Get second Repository",
			givenRepoName: "repo2",
			wantRepoName:  "repo2",
			wantInterface: (*scanItem.RepositoryInterface)(nil),
			wantErr:       nil,
		},
		{
			name:          "Get third Repository",
			givenRepoName: "repo3",
			wantRepoName:  "repo3",
			wantInterface: (*scanItem.RepositoryInterface)(nil),
			wantErr:       nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sut := NewRepositoryRegistry(memDb.NewMemDbConnection())
			got, err := sut.GetRepository(tt.givenRepoName)

			assert.Exactly(t, tt.wantErr, err)
			assert.Exactly(t, tt.wantRepoName, got.GetRepoName())
			assert.Implements(t, tt.wantInterface, got)
		})
	}
}

func TestRepositoryManager_GetDB__Given__Repo_Exist__Expect__Return_TheSameRepositoryObject(t *testing.T) {
	tests := []struct {
		name          string
		givenRepoName string
		wantRepoName  string
		wantObject    scanItem.RepositoryInterface
		wantErr       error
	}{
		{
			name:          "Get first Repository",
			givenRepoName: "repo1",
			wantRepoName:  "repo1",
			wantObject:    nil,
			wantErr:       nil,
		},
		{
			name:          "Get second Repository",
			givenRepoName: "repo2",
			wantRepoName:  "repo2",
			wantObject:    nil,
			wantErr:       nil,
		},
		{
			name:          "Get third Repository",
			givenRepoName: "repo3",
			wantRepoName:  "repo3",
			wantObject:    nil,
			wantErr:       nil,
		},
	}

	// ARRANGE
	sut := NewRepositoryRegistry(memDb.NewMemDbConnection())

	// call GetRepository the first time => repo object will be created
	for i, tt := range tests {
		tests[i].wantObject, _ = sut.GetRepository(tt.givenRepoName)
	}

	for _, tt := range tests {
		// ACT

		// call GetRepository() second time (with same repoName)
		// => the same Repository Object will be returned
		got, err := sut.GetRepository(tt.givenRepoName)

		// ASSERT
		assert.Exactly(t, tt.wantErr, err)
		assert.Exactly(t, tt.wantObject, got)
		assert.Exactly(t, tt.wantRepoName, got.GetRepoName())
	}
}