
## struct: FileSystem
- **CacheDirectory** - string
  - Path of cache directory [[Ref]](https://github.com/cs3org/cato/tree/master/examples/filesystem.go#L6)
  - Default: "/var/tmp/"
- **EnableLogging** - bool
  - Whether to enable logging [[Ref]](https://github.com/cs3org/cato/tree/master/examples/filesystem.go#L7)
  - Default: false
- **AvailableChecksums** - []string
  - The list of checksums provided by the file system [[Ref]](https://github.com/cs3org/cato/tree/master/examples/filesystem.go#L8)
  - Default: [adler, rabin]
- **DriverConfig** - map[string]map[string]interface{}
  - Configs for various metadata drivers [[Ref]](https://github.com/cs3org/cato/tree/master/examples/filesystem.go#L10)
  - Default: {json:{encoding: UTF8}, xml:{encoding: ASCII}}
- **Uploads** - *UploadConfig
  - Config for the HTTP uploads service [[Ref]](https://github.com/cs3org/cato/tree/master/examples/filesystem.go#L12)
  - Default: &UploadConfig{HTTPPrefix: uploads, DisableTus: false}

## struct: UploadConfig
- **disable_tus** - bool
  - Whether to disable TUS protocol for uploads. [[Ref]](https://github.com/cs3org/cato/tree/master/examples/filesystem.go#L17)
  - Default: false
- **http_prefix** - string
  - The prefix at which the uploads service should be exposed. [[Ref]](https://github.com/cs3org/cato/tree/master/examples/filesystem.go#L19)
  - Default: "uploads"
