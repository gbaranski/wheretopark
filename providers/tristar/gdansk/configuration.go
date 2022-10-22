package gdansk

import (
	"os"
	wheretopark "wheretopark/go"

	"gopkg.in/yaml.v2"
)

type Configuration struct {
	ParkingLots map[wheretopark.ID]wheretopark.Metadata `yaml:"parking-lots"`
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
