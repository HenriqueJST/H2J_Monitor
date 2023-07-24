package config

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
}

type monitor struct {
	LogsPath      string
	StagePath     string
	QlikLogPath   string
	Cliente       string
	EngineHost    string
	EnginePort    string
	UserName      string
	UserDirectory string
	Saas          string
}

// MonitorConfig instance of log settings
var MonitorConfig = &monitor{}

// Setup initialize the configuration instance

func Setup() {
	exePath, _ := os.Executable()
	dir := filepath.Dir(exePath)
	err := godotenv.Load(dir + "/.env")
	if err != nil {
		return
	}

	MonitorConfig.LogsPath = os.Getenv("LOGS_PATH")
	MonitorConfig.StagePath = os.Getenv("STAGE_PATH")
	MonitorConfig.Saas = os.Getenv("SAAS")
	if os.Getenv("SAAS") == "true" {
		MonitorConfig.QlikLogPath = os.Getenv("QLIK_LOGS_PATH_SAAS")
	} else {
		MonitorConfig.QlikLogPath = os.Getenv("QLIK_LOGS_PATH_ONPRIMESE")
	}
	MonitorConfig.Cliente = os.Getenv("CLIENTE_NAME")
	MonitorConfig.EngineHost = os.Getenv("ENGINEHOST")
	MonitorConfig.EnginePort = os.Getenv("ENGINEPORT")
	MonitorConfig.UserName = os.Getenv("USERNAME")
	MonitorConfig.UserDirectory = os.Getenv("USERDIRECTORY")
}
