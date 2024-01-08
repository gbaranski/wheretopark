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
	var internalConfiguration internalConfiguration
	err = json.Unmarshal(configurationJson, &internalConfiguration)
	if err != nil {
		return nil, err
	}

	configuration := Configuration{
		ParkingLots: make(map[string]ConfiguredParkingLot, len(internalConfiguration.ParkingLots)),
	}

	for _, parkingLot := range internalConfiguration.ParkingLots {
		id := wheretopark.GeometryToID(parkingLot.Geometry)
		totalSpots := make(map[wheretopark.SpotType]uint, len(parkingLot.Cameras))
		for _, camera := range parkingLot.Cameras {
			for _, spot := range camera.Spots {
				totalSpots[spot.Type]++
			}
		}
		parkingLot.TotalSpots = totalSpots
		configuration.ParkingLots[id] = parkingLot
	}

	return &configuration, nil
}

type internalConfiguration struct {
	ParkingLots []ConfiguredParkingLot `yaml:"parkingLots"`
}

type Configuration struct {
	ParkingLots map[wheretopark.ID]ConfiguredParkingLot `json:"parkingLots"`
}

type ConfiguredParkingLot struct {
	wheretopark.Metadata `json:",inline"`

	Cameras []Camera `json:"cameras"`
}

func (p *ConfiguredParkingLot) UnmarshalJSON(data []byte) error {
	var metadata wheretopark.Metadata
	err := json.Unmarshal(data, &metadata)
	if err != nil {
		return err
	}
	p.Metadata = metadata

	var cameras struct {
		Cameras []Camera `json:"cameras"`
	}
	err = json.Unmarshal(data, &cameras)
	if err != nil {
		return err
	}
	p.Cameras = cameras.Cameras
	return nil
}

type Camera struct {
	URL   string `json:"url"`
	Spots []Spot `json:"spots"`
}

type Spot struct {
	Points []Point              `json:"points"`
	Type   wheretopark.SpotType `json:"type,omitempty"`
}

func (s *Spot) UnmarshalJSON(data []byte) error {
	s.Type = wheretopark.SpotTypeCar // set default value before unmarshaling
	type Alias Spot                  // create alias to prevent endless loop
	return json.Unmarshal(data, (*Alias)(s))
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
