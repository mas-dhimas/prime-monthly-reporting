package config

import (
	"errors"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type ConfigLoader struct {
	file string
}

// NewConfig create new ConfigLoader to GetServiceConfig()
func NewConfig(file string) *ConfigLoader {
	return &ConfigLoader{
		file: file,
	}
}

func (c *ConfigLoader) GetServiceConfig() (*ServiceConfig, error) {
	config := &ServiceConfig{}

	err := cleanenv.ReadConfig(c.file, config)
	if err != nil {
		// Return early if the error is not 'file not found'
		if !errors.Is(err, os.ErrNotExist) {
			return nil, err
		}

		// If the error is 'file not found', try reading from environment variables
		err = cleanenv.ReadEnv(config)
		if err != nil {
			return nil, err
		}
	}
	return config, nil
}

// ServiceConfig stores the whole configuration for service.
type ServiceConfig struct {
	ServiceData ServiceDataConfig `yaml:"service_data"`
	SourceData  SourceDataConfig  `yaml:"source_data"`
}

// ServiceDataConfig contains the service data configuration.
type ServiceDataConfig struct {
	LogLevel   string `yaml:"log_level" env:"SERVICE_DATA_LOG_LEVEL"`
	PrimeToken string `yaml:"prime_token" env:"SERVICE_DATA_PRIME_TOKEN"`
}

// SourceDataConfig contains the source data configuration.
type SourceDataConfig struct {
	Network NetworkSourceData `yaml:"network"`
}

type NetworkSourceData struct {
	DeviceInventory    string `yaml:"device_inventory" env:"NETWORK_DEVICE_INVENTORY_SOURCE_DATA"`
	DeviceAvailability string `yaml:"device_availability" env:"NETWORK_DEVICE_AVAILABILITY_SOURCE_DATA"`
}
