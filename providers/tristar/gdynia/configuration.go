package gdynia

import (
	"os"
	wheretopark "wheretopark/go"

	_ "embed"

	"gopkg.in/yaml.v2"
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
	ParkingLots map[int]wheretopark.Metadata `json:"parkingLots"`
}

// Load Configuration from a YAML file
func LoadConfiguration(path string) (*Configuration, error) {
	var configuration Configuration
	yamlFile, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(yamlFile, &configuration)
	if err != nil {
		return nil, err
	}

	return &configuration, nil
}

func ParseConfiguration(content string) (*Configuration, error) {
	var configuration Configuration
	err := yaml.Unmarshal([]byte(content), &configuration)
	if err != nil {
		return nil, err
	}

	return &configuration, nil
}
