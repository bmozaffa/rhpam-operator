package defaults

//go:generate sh -c "CGO_ENABLED=0 go run .packr/packr.go $PWD"

import (
	"github.com/imdario/mergo"
	"gopkg.in/yaml.v2"

	"github.com/gobuffalo/packr"
)

func ConsoleEnvironmentDefaults() map[string]string {
	return overrideDefaults("console-env.yaml")
}

func ServerEnvironmentDefaults() map[string]string {
	return overrideDefaults("server-env.yaml")
}

func overrideDefaults(filename string) map[string]string {
	defaults := loadYamlMap("common-env.yaml")
	configuration := loadYamlMap(filename)
	mergo.Map(&defaults, configuration, mergo.WithOverride)
	return defaults
}

func loadYamlMap(filename string) map[string]string {
	box := packr.NewBox("../../../config/app")
	yamlMap := make(map[string]string)
	yaml.Unmarshal(box.Bytes(filename), &yamlMap)
	return yamlMap
}
