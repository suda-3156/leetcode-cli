package file

import (
	"bytes"
	"os"
	"text/template"
	"time"
)

const (
	BACKUP_TMPL      = "{{ .OriginalFile }}.bak.{{ .Timestamp }}"
	TIMESTAMP_FORMAT = "20060102150405"
)

func GetBackupPath(originalFile, timestamp string) string {
	tmpl, err := template.New("backupTemplate").Parse(BACKUP_TMPL)
	if err != nil {
		panic("failed to parse backup template: " + err.Error())
	}

	var buf bytes.Buffer
	data := map[string]string{
		"OriginalFile": originalFile,
		"Timestamp":    timestamp,
	}

	if err := tmpl.Execute(&buf, data); err != nil {
		panic("failed to execute backup template: " + err.Error())
	}
	return buf.String()
}

func CreateBackup(originalFilePath string) (string, error) {
	timestamp := time.Now().Format(TIMESTAMP_FORMAT)

	backupPath := GetBackupPath(originalFilePath, timestamp)

	input, err := os.ReadFile(originalFilePath)
	if err != nil {
		return "", err
	}

	if err := os.WriteFile(backupPath, input, 0600); err != nil {
		return "", err
	}

	return backupPath, nil
}
