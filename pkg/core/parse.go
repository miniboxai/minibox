package core

var defaultVersion = "minibox.ai/train.v1"

func parseVersion(m map[interface{}]interface{}) (interface{}, error) {
	var ver string
	if version, ok := m["version"].(string); ok {
		ver = version
	} else {
		ver = defaultVersion
	}

	return getConfigStruct(ver)
}
