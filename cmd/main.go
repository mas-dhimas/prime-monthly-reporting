package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/sirupsen/logrus"
	"gitlab.com/mas-dhimas/xlsx-prime-monthly-reporting/config"
	"gitlab.com/mas-dhimas/xlsx-prime-monthly-reporting/internal/network/handler"
	"gitlab.com/mas-dhimas/xlsx-prime-monthly-reporting/internal/network/repository"
	service "gitlab.com/mas-dhimas/xlsx-prime-monthly-reporting/internal/network/service"

	nhttp "gitlab.com/mas-dhimas/xlsx-prime-monthly-reporting/pkg/http"
)

const (
	DEVELOPER  = "Adhimas W Ramadhana"
	ASCIIImage = `
      ___           ___                         ___           ___     
     /\  \         /\__\                       /\  \         /\  \    
     \:\  \       /:/ _/_         ___         |::\  \       /::\  \   
      \:\  \     /:/ /\__\       /\__\        |:|:\  \     /:/\:\  \  
  _____\:\  \   /:/ /:/ _/_     /:/  /      __|:|\:\  \   /:/  \:\  \ 
 /::::::::\__\ /:/_/:/ /\__\   /:/__/      /::::|_\:\__\ /:/__/ \:\__\
 \:\~~\~~\/__/ \:\/:/ /:/  /  /::\  \      \:\~~\  \/__/ \:\  \ /:/  /
  \:\  \        \::/_/:/  /  /:/\:\  \      \:\  \        \:\  /:/  / 
   \:\  \        \:\/:/  /   \/__\:\  \      \:\  \        \:\/:/  /  
    \:\__\        \::/  /         \:\__\      \:\__\        \::/  /   
     \/__/         \/__/           \/__/       \/__/         \/__/    
	`
)

func main() {
	customFormatter := &logrus.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			return "", fmt.Sprintf(" %s:%d", f.File, f.Line)
		},
	}
	logrus.SetFormatter(customFormatter)
	logrus.SetReportCaller(true)
	logrus.SetLevel(logrus.InfoLevel)

	configPath := flag.String("config", "", "config file path")
	flag.Parse()

	// Get Config
	configLoader := config.NewConfig(*configPath)
	serviceConfig, err := configLoader.GetServiceConfig()
	if err != nil {
		logrus.Fatalf("Unable to load configuration: %v", err)
	}

	level, err := logrus.ParseLevel(serviceConfig.ServiceData.LogLevel)
	if err != nil {
		logrus.Fatalf("Unable to read log level : %s", err.Error())
	} else {
		logrus.SetLevel(level)
	}

	// Pre-printed text at startup.
	fmt.Println("Netmonk Prime XLSX Monthly Reporting Summary")
	fmt.Printf("Developed by %v.\n", DEVELOPER)
	fmt.Println(ASCIIImage)
	fmt.Println("⚠️  This programs is BETA version, PLEASE ONLY GENERATE UNDER 50 DEVICE for now ⚠️")
	logrus.Infof("Start service...")

	nr := repository.NewRepository(nhttp.NewHttpClient(serviceConfig.ServiceData.PrimeToken), serviceConfig.SourceData.Network)
	ns := service.NewService(&nr)
	nh := handler.NewHandler(&ns)

	// initiate docs folder for generated xlsx documents
	path := "docs"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.Mkdir(path, 0755)
		if err != nil {
			logrus.Error("unable to initiate docs directory, err: ", err)
			os.Exit(1)
		}
	}

	done := make(chan struct{})

	go func() {
		if err := nh.RouteService(); err != nil {
			logrus.Error(err)
			close(done)
			return
		}

		close(done)
	}()

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-quit:
		logrus.Info("received shutdown signal")
	case <-done:
		logrus.Info("route service finished, exiting")
	}
	logrus.Info("Shutting down server...")
	logrus.Info("Server exited gracefully")

}
