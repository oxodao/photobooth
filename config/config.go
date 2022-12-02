package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

var GET Config

type PhotoboothConfig struct {
	HardwareFlash      bool `yaml:"hardware_flash" json:"hardware_flash"`
	DefaultTimer       int  `yaml:"default_timer" json:"-"`
	UnattendedInterval int  `yaml:"unattended_interval" json:"-"`
}

type Config struct {
	Web struct {
		ListeningAddr string `yaml:"listening_addr"`
	} `yaml:"web"`

	DebugMode bool `yaml:"debug_mode"`

	RootPath    string `yaml:"root_path"`
	DefaultMode string `yaml:"default_mode"`

	Photobooth PhotoboothConfig `yaml:"photobooth"`
}

func (c *Config) GetImageFolder(eventId int64, unattended bool) (string, error) {
	subfolder := "pictures"
	if unattended {
		subfolder = "unattended"
	}

	path := filepath.Join(c.RootPath, "images", fmt.Sprintf("%v", eventId), subfolder)
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return "", err
	}
	return path, nil
}

func Load() error {
	cfg := Config{}

	configPath := os.Getenv("PHOTOMATON_CONFIG_PATH")
	if len(configPath) == 0 {
		configPath = "/etc/photomaton.yaml"
	}

	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return err
	}

	if cfg.Photobooth.DefaultTimer == 0 {
		cfg.Photobooth.DefaultTimer = 3
	}

	if len(cfg.DefaultMode) == 0 {
		cfg.DefaultMode = "PHOTOBOOTH"
	}

	GET = cfg

	return nil
}
