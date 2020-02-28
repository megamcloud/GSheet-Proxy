package config

type ConfigurationInfo struct {
	Server    string     `yaml:"server"`
	Monitor   string     `yaml:"monitor"`
	Storage   Storage    `yaml:"storage"`
	DbSources []DbSource `yaml:"dbsources"`
	Logging   Logger     `yaml:"logging"`
}

type Storage struct {
	Adapter string `yaml:"adapter"`
	Folder string `yaml:"folder"`
}

// DbSourceConfig ...
type DbSource struct {
	Name           string `yaml:"name"`
	IdField        string `yaml:"idfield"`
	FetchingUrl    string `yaml:"fetchingurl"`
	FetchingFormat string `yaml:"fetchingformat"`
	UpdateUrl      string `yaml:"updateurl"`
	UpdateMethod   string `yaml:"updatemethod"`
}

// LoggerConfig ....
type Logger struct {
	Filename   string `yaml:"filename"`
	MaxBackups int    `yaml:"MaxBackups"`
	MaxSize    int    `yaml:"MaxSize"` // megabytes
	MaxAge     int    `yaml:"MaxAge"`  // days
}
