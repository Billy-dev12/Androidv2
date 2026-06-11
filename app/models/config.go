package models

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type ConfigModel struct {
	DefaultDeviceID string `json:"default_device_id"`
	OutputDir       string `json:"output_dir"`
	configPath      string
	dataDir         string
}

func NewConfigModel() *ConfigModel {
	home, _ := os.UserHomeDir()
	dataDir := filepath.Join(home, ".android-tool")
	configPath := filepath.Join(dataDir, "config.json")
	return &ConfigModel{
		configPath: configPath,
		dataDir:    dataDir,
	}
}

func (c *ConfigModel) DataDir() string {
	return c.dataDir
}

func (c *ConfigModel) ConfigPath() string {
	return c.configPath
}

func (c *ConfigModel) Load() error {
	if err := c.ensureDir(); err != nil {
		return err
	}
	data, err := os.ReadFile(c.configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	return json.Unmarshal(data, c)
}

func (c *ConfigModel) Save() error {
	if err := c.ensureDir(); err != nil {
		return err
	}
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(c.configPath, data, 0644)
}

func (c *ConfigModel) Set(key, value string) error {
	switch key {
	case "default_device_id", "device_id":
		c.DefaultDeviceID = value
	case "output_dir":
		c.OutputDir = value
	default:
		return fmt.Errorf("unknown config key: %s (valid: default_device_id, output_dir)", key)
	}
	return c.Save()
}

func (c *ConfigModel) Show() string {
	var sb strings.Builder
	sb.WriteString("=== CONFIGURATION ===\n")
	sb.WriteString(fmt.Sprintf("  Config File : %s\n", c.configPath))
	sb.WriteString(fmt.Sprintf("  Default Device ID : %s\n", ifEmpty(c.DefaultDeviceID, "(not set)")))
	sb.WriteString(fmt.Sprintf("  Output Directory  : %s\n", ifEmpty(c.OutputDir, "(current dir)")))
	return sb.String()
}

func (c *ConfigModel) ensureDir() error {
	return os.MkdirAll(c.dataDir, 0755)
}

func ifEmpty(s, fallback string) string {
	if s == "" {
		return fallback
	}
	return s
}
