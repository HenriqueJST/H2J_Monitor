package monitor

import (
	"os"
	"time"

	"bitbucket.org/h2j/h2j-qlik-monitor/src/config"
	"bitbucket.org/h2j/h2j-qlik-monitor/src/qlik"
	"bitbucket.org/h2j/h2j-qlik-monitor/src/server"
	"github.com/fsnotify/fsnotify"
	"github.com/kardianos/service"
)

type Monitor struct {
	slog service.Logger
}

func (m *Monitor) Start(s service.Service) error {
	m.slog.Info("Started")
	go m.run()
	return nil
}
func (m *Monitor) run() {
	// watch qlik logs folder
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		m.slog.Error("NewWatcher failed: ", err)
	}
	defer watcher.Close()

	m.slog.Info("Adding watchers to ", config.MonitorConfig.QlikLogPath)
	err = addWatcher(config.MonitorConfig.QlikLogPath, watcher, m.slog)
	if err != nil {
		m.slog.Error("AddWatcher failed: ", err)
	}

	logsCheck := qlik.LogsCheck{}

	done := make(chan bool)
	go func() {
		defer close(done)

		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				err := logsCheck.CheckFile(event.Name, event.Op)
				if err != nil {
					m.slog.Error("file:", event.Name, " error: ", err)
					return
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				m.slog.Error("error:", err)
			}
		}

	}()

	go func() {
		for {
			server.Server_Service()
			time.Sleep(10 * time.Minute)
		}
	}()

	if config.MonitorConfig.Saas == "false" {
		go func() {
			for {
				var engineHost = config.MonitorConfig.EngineHost
				var enginePort = config.MonitorConfig.EnginePort
				var userName = config.MonitorConfig.UserName
				var userDirectory = config.MonitorConfig.UserDirectory
				qlik.Engine_Getlist(engineHost, enginePort, userName, userDirectory)
				time.Sleep(6 * time.Hour)
			}
		}()
	}

	<-done
}

func (m *Monitor) Stop(s service.Service) error {
	m.slog.Info("Stoped")
	return nil
}

func (m *Monitor) SetLogger(l service.Logger) {
	m.slog = l
}

func addWatcher(dir string, watcher *fsnotify.Watcher, logger service.Logger) (err error) {
	dirList, _ := os.ReadDir(dir)
	for _, item := range dirList {
		if item.IsDir() {
			dirName := dir + "\\" + item.Name()
			logger.Info("Add watcher to ", dirName)
			watcher.Add(dirName)
			addWatcher(dirName, watcher, logger)

		}
	}
	return
}
