package env

import (
	"fmt"
	"main/internal/cfg"
	"os"

	"github.com/sirupsen/logrus"
)

type Environment struct {
	C cfg.Config
}

var E *Environment

var (
	ConfigPath = "CONFIG_PATH"
)

func init() {
	// Get config path from environment variable
	path := os.Getenv(ConfigPath)
	if path == "" {
		path = "config.yaml"
	}

	var err error
	E, err = newEnvironment(path)
	if err != nil {
		logrus.Panic(fmt.Errorf("failed to load config: %w", err))
	}

	configureLogger()
}

func newEnvironment(yamlFile string) (*Environment, error) {
	conf, err := cfg.NewConfig(yamlFile)
	if err != nil {
		return nil, err
	}

	return &Environment{C: conf}, nil
}

func configureLogger() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	if E.C.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
}
