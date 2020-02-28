package config

import (
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"os"
	"testing"
	"time"
)

func TestConfiguration_GetConfig__Given__NoFolderExist_And_NoConfigFileExist__Expect__Panic(t *testing.T) {
	// arrange mock file system
	mockFs := afero.NewMemMapFs()
	sut := NewConfig().mockFileSystem(mockFs)

	assert.Panics(t, func() { sut.GetConfig() })
}

func TestConfiguration_GetConfig__Given__AllFoldersExist_CorrectPermission_And_ConfigFileExist__Expect__NotPanic(t *testing.T) {

	mockFs := afero.NewMemMapFs()
	sut := NewConfig().mockFileSystem(mockFs)

	// arrange mock file system
	setupMockFolder(t, sut, mockFs, 0755)
	configFile := setupGoodMockConfigFile(t, sut, mockFs, 0600)

	assert.NotPanics(t, func() { sut.GetConfig() })
	assert.Exactly(t, configFile, sut.getUsedConfigFile())
}

func TestConfiguration_GetConfig__Given__AllFoldersExist_CorrectPermission_And_BadConfigFileExist__Expect__Panic(t *testing.T) {

	mockFs := afero.NewMemMapFs()
	sut := NewConfig().mockFileSystem(mockFs)

	// arrange mock file system
	setupMockFolder(t, sut, mockFs, 0755)
	setupBadMockConfigFile(t, sut, mockFs, 0600)

	assert.Panics(t, func() { sut.GetConfig() })
}

func TestConfiguration_GetConfig__Given__AllFoldersExist_CorrectPermission_And_GoodConfigFileExist__Expect__ReturnStructure(t *testing.T) {

	mockFs := afero.NewMemMapFs()
	sut := NewConfig().mockFileSystem(mockFs)

	// arrange mock file system
	setupMockFolder(t, sut, mockFs, 0755)
	setupGoodMockConfigFile(t, sut, mockFs, 0600)

	assert.Exactly(t, goodMockConfigStruct1(), sut.GetConfig())
}

func setupMockFolder(t *testing.T, sut *Configuration, mockFs afero.Fs, perm os.FileMode) {
	folders := sut.getLookupFolders()

	for _, folder := range folders {
		if err := mockFs.MkdirAll(folder, perm); err != nil {
			panic(err.Error())
		}

		exist, _ := afero.DirExists(mockFs, folder)
		assert.True(t, exist)
	}
}

func setupGoodMockConfigFile(t *testing.T, sut *Configuration, mockFs afero.Fs, perm os.FileMode) string {

	rand.Seed(time.Now().UnixNano())
	folders := sut.getLookupFolders()
	randFolder := folders[rand.Intn(len(folders)-1)]

	configFile := randFolder + "/" + sut.getConfigName() + ".yaml"
	_ = afero.WriteFile(mockFs, configFile, []byte(goodMockConfigYAMLContent1()), perm)

	exist, _ := afero.Exists(mockFs, configFile)
	assert.True(t, exist)

	return configFile
}

func setupBadMockConfigFile(t *testing.T, sut *Configuration, mockFs afero.Fs, perm os.FileMode) string {

	rand.Seed(time.Now().UnixNano())
	folders := sut.getLookupFolders()
	randFolder := folders[rand.Intn(len(folders)-1)]

	configFile := randFolder + "/" + sut.getConfigName() + ".yaml"
	_ = afero.WriteFile(mockFs, configFile, []byte(badMockConfigContent()), perm)

	exist, _ := afero.Exists(mockFs, configFile)
	assert.True(t, exist)

	return configFile
}


func goodMockConfigStruct1() *ConfigurationInfo{
	return &ConfigurationInfo{
		Server:    "0.0.0.0:443",
		Monitor:   "0.0.0.0:2222",
		DbSources: []DbSource{
			{
				Name: "name_0",
				FetchingUrl: "fetching_url_0",
				FetchingFormat: "json",
				UpdateUrl: "update_url_0",
				UpdateMethod: "GET",
			},
			{
				Name: "name_1",
				FetchingUrl: "fetching_url_1",
				FetchingFormat: "json",
				UpdateUrl: "update_url_1",
				UpdateMethod: "GET",
			},
		},
		Logging: Logger{
			Filename:   "/var/log/event-hub/event-hub.log",
			MaxSize:    100,
			MaxBackups: 3,
			MaxAge:     28,
		},
	}
}

func goodMockConfigYAMLContent1() string {
	return `
server: 0.0.0.0:443
monitor: 0.0.0.0:2222

dbsources:
  - Name: name_0
    FetchingUrl: fetching_url_0
    FetchingFormat: json
    UpdateUrl: update_url_0
    UpdateMethod: GET

  - Name: name_1
    FetchingUrl: fetching_url_1
    FetchingFormat: json
    UpdateUrl: update_url_1
    UpdateMethod: GET

logging:
  Filename: /var/log/event-hub/event-hub.log
  MaxSize: 100
  MaxBackups: 3
  MaxAge: 28
`
}

func badMockConfigContent() string {
	return `
: 0.0.0.0:443
monitor: 0.0.0.0:2222

dbsources:
  abc:
    fetching_url: https://script.google.com/macros/s/AKfycbyKxlzZMiVlF01ZGPAXYsY0ARV-L8V04QCgONo5kIbTAwkfOC4C/exec?path=/sample&order=field_qrcode&limit=25&offset=:offset
    fetching_format: json
    update_url: https://script.google.com/macros/s/AKfycbyKxlzZMiVlF01ZGPAXYsY0ARV-L8V04QCgONo5kIbTAwkfOC4C/exec?path=/sample/:id&method=POST&field_khungdien=xyz&field_Salutation=caca5
    update_method: GET

  def:
    fetching_url: https://script.google.com/macros/s/AKfycbyKxlzZMiVlF01ZGPAXYsY0ARV-L8V04QCgONo5kIbTAwkfOC4C/exec?path=/sample&order=field_qrcode&limit=25&offset=:offset
    fetching_format: json
    update_url: https://script.google.com/macros/s/AKfycbyKxlzZMiVlF01ZGPAXYsY0ARV-L8V04QCgONo5kIbTAwkfOC4C/exec?path=/sample/:id&method=POST&field_khungdien=xyz&field_Salutation=caca5
    update_method: GET

logging:
  Filename: ./go-qcoordinator.log
  MaxSize: 100
  MaxBackups: 3
  MaxAge: 28
`
}
