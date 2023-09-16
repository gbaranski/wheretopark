package cctv

import (
	"encoding/json"
	"fmt"
	"image"
	"os"
	wheretopark "wheretopark/go"

	_ "embed"

	"github.com/goccy/go-yaml"
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
	yamlFile, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return ParseConfiguration(string(yamlFile))
}

func ParseConfiguration(content string) (*Configuration, error) {
	var configurationAny map[string]any
	err := yaml.Unmarshal([]byte(content), &configurationAny)
	if err != nil {
		return nil, err
	}
	configurationJson, err := json.Marshal(configurationAny)
	if err != nil {
		return nil, err
	}
	var configuration Configuration
	err = json.Unmarshal(configurationJson, &configuration)
	if err != nil {
		return nil, err
	}

	for i := range configuration.ParkingLots {
		parkingLot := &configuration.ParkingLots[i]
		parkingLot.TotalSpots = make(map[string]uint)
		totalSpots := 0
		for _, camera := range parkingLot.Cameras {
			totalSpots += len(camera.Spots)
		}

		parkingLot.TotalSpots["CAR"] = uint(totalSpots)
	}

	return &configuration, nil
}

type Configuration struct {
	ParkingLots []CameraParkingLot `json:"parkingLots"`
}

type CameraParkingLot struct {
	wheretopark.Metadata `json:",inline"`

	Cameras []ParkingLotCamera `json:"cameras"`
}

func (p *CameraParkingLot) UnmarshalJSON(data []byte) error {
	var metadata wheretopark.Metadata
	err := json.Unmarshal(data, &metadata)
	if err != nil {
		return err
	}
	p.Metadata = metadata

	var cameras struct {
		Cameras []ParkingLotCamera `json:"cameras"`
	}
	err = json.Unmarshal(data, &cameras)
	if err != nil {
		return err
	}
	p.Cameras = cameras.Cameras
	return nil
}

type ParkingLotCamera struct {
	URL   string        `json:"url"`
	Spots []ParkingSpot `json:"spots"`
}

type ParkingSpot struct {
	Points []Point `json:"points"`
}

type Point struct {
	X int
	Y int
}

func (p *Point) ToImagePoint() image.Point {
	return image.Point{X: p.X, Y: p.Y}
}

func (p *Point) UnmarshalJSON(b []byte) error {
	var array []int
	if err := json.Unmarshal(b, &array); err != nil {
		return err
	}
	if len(array) != 2 {
		return fmt.Errorf("expected array of length 2, got %d", len(array))
	}
	p.X = array[0]
	p.Y = array[1]
	return nil
}

func (p Point) MarshalJSON() ([]byte, error) {
	array := []int{p.X, p.Y}
	return json.Marshal(array)
}
