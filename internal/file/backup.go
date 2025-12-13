package file

import (
	"bytes"
	"fmt"
	"os"
	"text/template"
	"time"
)

const (
	BACKUP_TMPL      = "{{ .OriginalFile }}.bak.{{ .Timestamp }}"
	TIMESTAMP_FORMAT = "20060102150405"
)

func GetBackupPath(originalFile, timestamp string) (string, error) {
	tmpl, err := template.New("backupTemplate").Parse(BACKUP_TMPL)
	if err != nil {
		return "", fmt.Errorf("failed to parse backup template: %w", err)
	}

	var buf bytes.Buffer
	data := map[string]string{
		"OriginalFile": originalFile,
		"Timestamp":    timestamp,
	}

	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute backup template: %w", err)
	}
	return buf.String(), nil
}

func CreateBackup(originalFilePath string) (string, error) {
	timestamp := time.Now().Format(TIMESTAMP_FORMAT)

	backupPath, err := GetBackupPath(originalFilePath, timestamp)
	if err != nil {
		return "", err
	}

	input, err := os.ReadFile(originalFilePath)
	if err != nil {
		return "", err
	}

	if err := os.WriteFile(backupPath, input, 0600); err != nil {
		return "", err
	}

	return backupPath, nil
}
