package main

import (
	"flag"
	"os"

	bosherr "github.com/cloudfoundry/bosh-agent/errors"
	boshlog "github.com/cloudfoundry/bosh-agent/logger"
	boshsys "github.com/cloudfoundry/bosh-agent/system"

	"github.com/frodenas/bosh-google-cpi/registry/server"
	"github.com/frodenas/bosh-google-cpi/registry/server/store"
)

const mainLogTag = "main"

var (
	configPathOpt = flag.String("configPath", "", "Path to configuration file")
)

func main() {
	logger := boshlog.NewWriterLogger(boshlog.LevelDebug, os.Stderr, os.Stderr)
	fs := boshsys.NewOsFileSystem(logger)

	defer logger.HandlePanic("Main")

	flag.Parse()

	config, err := NewConfigFromPath(*configPathOpt, fs)
	if err != nil {
		logger.Error(mainLogTag, "Loading config: %s", err.Error())
		os.Exit(1)
	}

	instanceHandler, err := createInstanceHandler(config, logger)
	if err != nil {
		logger.Error(mainLogTag, "Creating an Instance Handler: %s", err.Error())
		os.Exit(1)
	}

	listener := server.NewListener(config.Server, instanceHandler, logger)
	err = listener.ListenAndServe()
	if err != nil {
		logger.Error(mainLogTag, "Starting Server: %s", err.Error())
		os.Exit(1)
	}
	listener.WaitForServerToExit()
}

func createInstanceHandler(config Config, logger boshlog.Logger) (*server.InstanceHandler, error) {
	registryStore, err := store.NewRegistryStore(config.Store, logger)
	if err != nil {
		return nil, bosherr.WrapError(err, "Creating a Registry Store")
	}

	instanceHandler := server.NewInstanceHandler(config.Server, registryStore, logger)

	return instanceHandler, nil
}
