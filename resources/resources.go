package resources

type DocumentationInfo struct {
	FieldName    string
	DataType     string
	DefaultValue string
	Description  string
}

type CatoConfig struct {
	CustomTag    string
	Driver       string
	DriverConfig map[string]map[string]interface{}
}
