package config

import (
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/afero"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
)

// NOTE: this must be not include extension such as .yaml, .json .....
const FileName = "config"
const AppName = "event-hub"

// Configuration structure
type ConfigurationInterface interface {
	GetConfig() *ConfigurationInfo
}

type Configuration struct {
	fileName      string
	lookupFolders []string
	viper         *viper.Viper
	content       *ConfigurationInfo
}

func NewConfig() *Configuration {
	c := &Configuration{
		fileName: FileName,
		viper:    viper.New(),
		content:  nil,
	}

	c.setupLookupFolders()

	return c
}

func (c *Configuration) setupLookupFolders() {
	//  "/etc/AppName"
	var folders = []string{"/etc/" + AppName}

	//   "/currentFolder"
	if dir, err := os.Getwd(); err == nil {
		folders = append(folders, dir)
	}

	//  "$HOME/.AppName"
	if dir, err := homedir.Dir(); err == nil {
		folders = append(folders, dir+"/."+AppName)
	}

	c.lookupFolders = folders
}

func (c *Configuration) mockFileSystem(fs afero.Fs) *Configuration {
	c.viper.SetFs(fs)
	return c
}

func (c *Configuration) GetConfig() *ConfigurationInfo {
	if nil == c.content {
		c.content = &ConfigurationInfo{}
	}

	if err := c.parseConfig(); err != nil {
		panic(err.Error())
	}

	return c.content
}


// ParseConfig will find and Parse Config
func (c *Configuration) parseConfig() error {
	if err := c.readConfig(); err != nil {
		return err
	}

	return c.viper.Unmarshal(&(c.content))
}

// ParseConfig will find and Parse Config
func (c *Configuration) readConfig() error {
	pflag.CommandLine = pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)

	filepath := pflag.StringP("config", "c", "na", "the path of the configuration file")
	pflag.Parse()

	if (*filepath) != "na" {
		c.viper.SetConfigFile(*filepath)
	} else {
		c.viper.SetConfigName(c.fileName)

		for _, folder := range c.lookupFolders {
			c.viper.AddConfigPath(folder)
		}
	}

	// Find and read the config file
	err := c.viper.ReadInConfig()

	return err
}

//func (c *Configuration) watchConfig() {
//	c.viper.WatchConfig()
//	c.viper.OnConfigChange(func(e fsnotify.Event) {
//	 	fmt.Println("Config file changed:", e.Name)
//	})
//}

// for Testing Only

func (c *Configuration) getConfigName() string {
	return c.fileName
}

func (c *Configuration) getLookupFolders() []string {
	return c.lookupFolders
}

func (c *Configuration) getUsedConfigFile() string {
	return c.viper.ConfigFileUsed()
}