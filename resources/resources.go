package resources

type FieldInfo struct {
	FieldName    string
	DataType     string
	DefaultValue string
	Description  string
	LineNumber   int
}

type CatoConfig struct {
	CustomTag    string
	Driver       string
	DriverConfig map[string]map[string]interface{}
}
