package env

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

const DefaultRootPath = "."

// util key
const (
	ContextKeyRequestID = "requestId"
	ContextKeyLogID     = "logID"
	ContextKeyNoLog     = "_no_log"
	ContextKeyUri       = "_uri"
	zapLoggerAddr       = "_zap_addr"
	sugaredLoggerAddr   = "_sugared_addr"
	customerFieldKey    = "__customerFields"
)

var (
	LocalIP string
	AppName string
	RunMode string

	Namespace   string
	ServiceName string

	runEnv int

	rootPath        string
	dockerPlateForm bool
)

// IsDockerPlatform 判断项目运行平台：容器 vs 开发环境
func IsDockerPlatform() bool {
	return dockerPlateForm
}

// SetAppName 开发环境可手动指定SetAppName
func SetAppName(appName string) {
	if !dockerPlateForm {
		AppName = appName
	}
}

func GetAppName() string {
	return AppName
}

// SetRootPath 设置应用的根目录
func SetRootPath(r string) {
	if !dockerPlateForm {
		rootPath = r
	}
}

// GetRootPath 返回应用的根目录
func GetRootPath() string {
	if rootPath != "" {
		return rootPath
	} else {
		return DefaultRootPath
	}
}

// GetConfDirPath 返回配置文件目录绝对地址
func GetConfDirPath() string {
	return filepath.Join(GetRootPath(), "conf")
}

// GetLogDirPath 返回log目录的绝对地址
func GetLogDirPath() string {
	return filepath.Join(GetRootPath(), "zap")
}

func GetRunEnv() int {
	return runEnv
}

const (
	SubConfDefault = ""
	SubConfMount   = "mount"
	SubConfApp     = "app"
)

func LoadConf(filename, subConf string, s interface{}) {
	var path string
	path = filepath.Join(GetConfDirPath(), subConf, filename)

	if yamlFile, err := os.ReadFile(path); err != nil {
		panic(filename + " get error: " + err.Error())
	} else if err = yaml.Unmarshal(yamlFile, s); err != nil {
		panic(filename + " unmarshal error: " + err.Error())
	}
}
