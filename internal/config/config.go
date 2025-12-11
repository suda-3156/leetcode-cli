package config

import "time"

const DefaultSearchLimit = 20

// GetDefaultOutputPath generates the default output path: ./YYYY-MM-DD/<frontendID>.<titleSlug><extension>
func GetDefaultOutputPath(frontendID, titleSlug, extension string) string {
	date := time.Now().Format("2006-01-02")
	return "./" + date + "/" + frontendID + "." + titleSlug + extension
}

func GetCurrentDate() string {
	return time.Now().Format("2006-01-02")
}
