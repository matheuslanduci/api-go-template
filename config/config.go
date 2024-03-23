package config

import (
	"fmt"
	"log"
	"os"
	"slices"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Mode   string `yaml:"mode"`
	Server struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	} `yaml:"server"`
	Cors struct {
		AllowedOrigins string `yaml:"allowed_origins"`
		AllowedMethods string `yaml:"allowed_methods"`
		AllowedHeaders string `yaml:"allowed_headers"`
	} `yaml:"cors"`
	Database struct {
		Url string `yaml:"url"`
	} `yaml:"database"`
}

func New(mode string) *Config {
	supportedModes := []string{"development", "production", "testing"}

	validMode := slices.Contains(supportedModes, mode)

	if !validMode {
		log.Fatalf("The current mode is not supported. Please use one of the following: %v", strings.Join(supportedModes, ", "))
	}

	path := fmt.Sprintf("%s.config.yaml", mode)

	fileInfo, err := os.Stat(path)

	if err != nil || fileInfo.IsDir() {
		log.Fatalf("A configuration file wasn't found for the current mode. Check if the file exists.")
	}

	config := &Config{}

	file, err := os.Open(path)

	if err != nil {
		log.Fatalf("An error occurred during the opening of the configuration file. Check your user's permissions.")
	}

	defer file.Close()

	yamlDecoder := yaml.NewDecoder(file)

	if err := yamlDecoder.Decode(&config); err != nil {
		log.Fatalf("Your configuration file is malstructured.")
	}

	return config
}
