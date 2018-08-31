package defaults

import (
	"encoding/json"

	"github.com/bmozaffa/rhpam-operator/config"
)

func ConsoleEnvironmentDefaults() map[string]string {
	return overrideDefaults("console-env.json")
}

func ServerEnvironmentDefaults() map[string]string {
	return overrideDefaults("server-env.json")
}

func overrideDefaults(filename string) map[string]string {
	defaults := loadJsonMap("common-env.json")
	configuration := loadJsonMap(filename)
	for key, value := range configuration {
		defaults[key] = value
	}
	return defaults
}

func loadJsonMap(filename string) map[string]string {
	box := config.ConfigPackr()
	jsonMap := make(map[string]string)
	json.Unmarshal(box.Bytes(filename), &jsonMap)
	return jsonMap
}
