// main.go
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
	logger := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			logger.Fatal(fmt.Errorf("Fatal error config file: %w \n", err))
		}
	}
	config := configuration.Configuration{}
	viper.Unmarshal(&config)
	listenPort := config.Server.ListenPort

	forwardUrl := config.Server.ForwardUrl

	logger.Printf("Listen port configured : %d\n", listenPort)
	logger.Printf("Forward url configured : %s\n", forwardUrl)
	proxy := reverse_proxy.CreateProxy(
		forwardUrl,
		uint32(listenPort),
		logger,
	)
	if config.HasMiddleware("log_request") {
		proxy.AppendHttpMiddewareAdapter(middleware.LogRequest(logger))
	}
	if config.HasMiddleware("dump_request") {
		proxy.AppendHttpMiddewareAdapter(middleware.DumpRequest(logger))
	}
	if config.HasMiddleware("cors") {
		proxy.AppendHttpResponseMiddewareAdapter(middleware.Cors())
	}

	proxy.Start()
}
