package registry

import "github.com/cs3org/cato/writer"

// NewFunc is the function prototype that drivers should register at init.
type NewFunc func(map[string]interface{}) (writer.ConfigWriter, error)

// NewFuncs is a map containing all the registered drivers.
var NewFuncs = map[string]NewFunc{}

// Register registers a new driver.
// Not safe for concurrent use, safe for use from package init.
func Register(name string, f NewFunc) {
	NewFuncs[name] = f
}
