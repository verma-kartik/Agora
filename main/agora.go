package main

import (
	"flag"
	log "github.com/cihub/seelog"
	"github.com/pkg/profile"
	"github.com/verma-kartik/Agora"
	"runtime"
)

func main() {
	// Set up a done channel, that's shared by the whole pipeline.
	// Closing this channel will kill all pipeline goroutines
	//done := make(chan struct{})
	//defer close(done)

	// Set up logging
	initializeLogging()

	//flush log before we shut down
	defer log.Flush()

	config := parseCommandLineFlags()
	Agora.SetConfig(&config)

	if config.ProfilingEnabled {
		defer profile.Start(profile.CPUProfile).Stop()
	}

	log.Infof("Broker started on port: %d", Agora.Configuration.Port)
	log.Infof("Executing on: %d threads", runtime.GOMAXPROCS(-1))

	connectionManager := Agora.NewConnectionManager()
	connectionManager.Start()
}
func initializeLogging() {
	logger, err := log.LoggerFromConfigAsFile("config/logconfig.xml")

	if err != nil {
		log.Criticalf("An error occurred whilst initializing logging\n", err.Error())
		panic(err)
	}

	log.ReplaceLogger(logger)
}

func parseCommandLineFlags() Agora.Config {
	configToReturn := Agora.Config{}

	flag.IntVar(&configToReturn.Port, "port", 48879, "The port to listen on")
	flag.BoolVar(&configToReturn.ProfilingEnabled, "profile", false, "Produce a pprof file")
	flag.StringVar(&configToReturn.StatsDEndpoint, "statsd", "", "The StatsD endpoint to send metrics to")

	flag.Parse()

	return configToReturn
}
