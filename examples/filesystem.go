package main

import "fmt"

type FileSystem struct {
	CacheDirectory     string   `docs:"/var/tmp/;Path of cache directory"`
	EnableLogging      bool     `docs:"false;Whether to enable logging"`
	AvailableChecksums []string `docs:"[adler, rabin];The list of checksums provided by the file system"`
	// Configs for various metadata drivers
	DriverConfig map[string]map[string]interface{} `docs:"{json:{encoding: UTF8}, xml:{encoding: ASCII}}"`
	// Config for the HTTP uploads service
	Uploads *UploadConfig `docs:"&UploadConfig{HTTPPrefix: uploads, DisableTus: false}"`
}

type UploadConfig struct {
	// Whether to disable TUS protocol for uploads.
	DisableTus bool `json:"disable_tus" docs:"false"`
	// The prefix at which the uploads service should be exposed.
	HTTPPrefix string `json:"http_prefix" docs:"uploads"`
}

func (fs FileSystem) init() {
	if fs.CacheDirectory == "" {
		fs.CacheDirectory = "/var/tmp/"
	}
	if len(fs.AvailableChecksums) == 0 {
		fs.AvailableChecksums = []string{"adler", "rabin"}
	}
	if fs.DriverConfig == nil {
		fs.DriverConfig = map[string]map[string]interface{}{
			"json": map[string]interface{}{
				"encoding": "UTF8",
			},
			"xml": map[string]interface{}{
				"encoding": "ASCII",
			},
		}
	}
	if fs.Uploads == nil {
		fs.Uploads = &UploadConfig{
			HTTPPrefix: "uploads",
		}
	}
}

func main() {
	fs := FileSystem{
		AvailableChecksums: []string{"adler32"},
	}
	fmt.Printf("FileSystem: %+v", fs)
}
