package cctv

import (
	"fmt"
	"image"
	"os"
	wheretopark "wheretopark/go"

	_ "embed"

	"gopkg.in/yaml.v3"
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

type Configuration struct {
	ParkingLots []ParkingLot `yaml:"parking-lots"`
}

type ParkingLot struct {
	wheretopark.Metadata `yaml:",inline"`

	CameraURL string        `yaml:"camera-url"`
	Spots     []ParkingSpot `yaml:"spots"`
}

type ParkingSpot struct {
	Points []Point `yaml:"points"`
}

type Point struct {
	X int
	Y int
}

func (p *Point) ToImagePoint() image.Point {
	return image.Point{X: p.X, Y: p.Y}
}

func (p *Point) UnmarshalYAML(node *yaml.Node) error {
	var array []int
	if err := node.Decode(&array); err != nil {
		return err
	}
	if len(array) != 2 {
		return fmt.Errorf("expected array of length 2, got %d", len(array))
	}
	p.X = array[0]
	p.Y = array[1]
	return nil
}
