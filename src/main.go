package main

import (
	"flag"
	"log"

	"bitbucket.org/h2j/h2j-qlik-monitor/src/config"
	"bitbucket.org/h2j/h2j-qlik-monitor/src/directory_verify"
	"bitbucket.org/h2j/h2j-qlik-monitor/src/monitor"
	"github.com/kardianos/service"
)

const serviceName = "H2JMonitor"
const serviceDisplayName = "H2J Monitor"
const serviceDescription = "Monitoramento dos serviços"

var (
	mode string
	slog service.Logger
)

func init() {
	config.Setup()
	flag.StringVar(&mode, "mode", "", "install/uninstall/run")
	flag.Parse()
}

func main() {
	svcConfig := &service.Config{
		Name:        serviceName,
		DisplayName: serviceDisplayName,
		Description: serviceDescription,
	}

	prg := &monitor.Monitor{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}
	errs := make(chan error, 5)
	slog, err = s.Logger(errs)
	if err != nil {
		log.Fatal(err)
	}
	prg.SetLogger(slog)

	go func() {
		for {
			err := <-errs
			if err != nil {
				log.Println(err)
			}
		}
	}()

	if mode == "install" {
		directory_verify.Verify_directory()

		err = s.Install()
		if err != nil {
			slog.Error(err)
		}
		slog.Info("Successfully installed")
	}

	if mode == "uninstall" {
		directory_verify.Delete_directory()
		err = s.Uninstall()
		if err != nil {
			slog.Error(err)
		}
		slog.Info("Successfully uninstalled")
	}

	if mode == "manual_uninstall" {
		err = s.Uninstall()
		if err != nil {
			slog.Error(err)
		}
		slog.Info("Successfully uninstalled")
	}

	if mode == "" || mode == "run" {
		directory_verify.Verify_directory()
		err = s.Run()
		if err != nil {
			slog.Error(err)
		}
	}
}
