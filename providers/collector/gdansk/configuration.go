package gdansk

import (
	"os"
	wheretopark "wheretopark/go"

	_ "embed"

	"github.com/ghodss/yaml"
)

//go:embed default.yaml
var defaultConfigurationContent string
var DefaultConfiguration Configuration

func init() {
	defaultConfiguration, err := ParseConfiguration(defaultConfigurationContent)
	if err != nil {
		panic(err)
	}
	DefaultConfiguration = *defaultConfiguration
}

type Configuration struct {
	ParkingLots map[string]wheretopark.Metadata `json:"parkingLots"`
}

// Load Configuration from a YAML file
func LoadConfiguration(path string) (*Configuration, error) {
	yamlFile, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return ParseConfiguration(string(yamlFile))
}

func ParseConfiguration(content string) (*Configuration, error) {
	var configuration Configuration
	err := yaml.Unmarshal([]byte(content), &configuration)
	if err != nil {
		return nil, err
	}

	return &configuration, nil
}
