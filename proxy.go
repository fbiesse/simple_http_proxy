package main

import (
	"fmt"
	"log"
	"os"

	"github.com/fbiesse/simple_reverse_proxy/configuration"
	"github.com/fbiesse/simple_reverse_proxy/reverse_proxy"
	"github.com/fbiesse/simple_reverse_proxy/reverse_proxy/middleware"

	"github.com/spf13/viper"
)

func main() {
	logflags := log.LstdFlags | log.Lshortfile
	stdoutLogger := log.New(os.Stdout, "", logflags)

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			stdoutLogger.Fatal(fmt.Errorf("Fatal error config file: %w \n", err))
		}
	}
	config := configuration.Configuration{}
	viper.Unmarshal(&config)
	listenPort := config.Server.ListenPort
	forwardUrl := config.Server.ForwardUrl

	stdoutLogger.Printf("Listen port configured : %d\n", listenPort)
	stdoutLogger.Printf("Forward url configured : %s\n", forwardUrl)

	var logger log.Logger

	if config.LogFile != "" {
		file, err := os.OpenFile(config.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			stdoutLogger.Fatal(err)
		}
		logger = *log.New(file, "", logflags)
		stdoutLogger.Printf("Logs configured to %s\n", config.LogFile)

	} else {
		logger = *stdoutLogger
		stdoutLogger.Printf("Logs configured to %s\n", "stdout")
	}

	proxy := reverse_proxy.CreateProxy(
		forwardUrl,
		uint32(listenPort),
		stdoutLogger,
	)
	if config.HasMiddleware("log_request") {
		proxy.AppendHttpMiddewareAdapter(middleware.LogRequest(&logger))
		stdoutLogger.Println("log_request middleware enabled")
	}
	if config.HasMiddleware("dump_request") {
		proxy.AppendHttpMiddewareAdapter(middleware.DumpRequest(&logger))
		stdoutLogger.Println("dump_request middleware enabled")
	}
	if config.HasMiddleware("cors") {
		proxy.AppendHttpResponseMiddewareAdapter(middleware.Cors())
		stdoutLogger.Println("cors middleware enabled")
	}

	proxy.Start()
}
