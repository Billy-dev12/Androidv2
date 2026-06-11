package models

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type HistoryModel struct {
	filePath string
	dataDir  string
}

func NewHistoryModel() *HistoryModel {
	home, _ := os.UserHomeDir()
	dataDir := filepath.Join(home, ".android-tool")
	historyPath := filepath.Join(dataDir, "history.log")
	return &HistoryModel{
		filePath: historyPath,
		dataDir:  dataDir,
	}
}

func (h *HistoryModel) FilePath() string {
	return h.filePath
}

func (h *HistoryModel) DataDir() string {
	return h.dataDir
}

func (h *HistoryModel) Append(command string, args ...string) error {
	if err := os.MkdirAll(h.dataDir, 0755); err != nil {
		return err
	}
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	allArgs := strings.Join(args, " ")
	line := fmt.Sprintf("[%s] %s %s\n", timestamp, command, allArgs)

	f, err := os.OpenFile(h.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(line)
	return err
}

func (h *HistoryModel) Show() (string, error) {
	data, err := os.ReadFile(h.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return "No history entries yet.", nil
		}
		return "", err
	}
	content := strings.TrimSpace(string(data))
	if content == "" {
		return "No history entries yet.", nil
	}
	var sb strings.Builder
	sb.WriteString("=== COMMAND HISTORY ===\n")
	sb.WriteString(fmt.Sprintf("  File: %s\n\n", h.filePath))
	sb.WriteString(content)
	sb.WriteString("\n")
	return sb.String(), nil
}

func (h *HistoryModel) Clear() error {
	return os.Remove(h.filePath)
}
