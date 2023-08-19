package gdansk

import (
	"time"

	"github.com/go-resty/resty/v2"
)

type Coordinate struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Metadata struct {
	ID             string     `json:"id"`
	Name           string     `json:"name"`
	ShortName      string     `json:"shortName"`
	Address        string     `json:"address"`
	StreetEntrance string     `json:"streetEntrance"`
	Location       Coordinate `json:"location"`
}

type State struct {
	ID             string `json:"parkingId"`
	AvailableSpots uint   `json:"availableSpots"`
	LastUpdate     string `json:"lastUpdate"`
}

type Response[T any] struct {
	LastUpdate  string `json:"lastUpdate"`
	ParkingLots []T    `json:"parkingLots"`
}

const (
	METADATA_URL = "https://ckan.multimediagdansk.pl/dataset/cb1e2708-aec1-4b21-9c8c-db2626ae31a6/resource/d361dff3-202b-402d-92a5-445d8ba6fd7f/download/parking-lots.jso"
	STATE_URL    = "https://ckan2.multimediagdansk.pl/parkingLots"
)

var client = resty.New()

func init() {
	client.GetClient().Timeout = 10 * time.Second
}

func GetMetadata() (*Response[Metadata], error) {
	resp, err := client.R().SetResult(&Response[Metadata]{}).Get(METADATA_URL)
	if err != nil {
		return nil, err
	}
	response := resp.Result().(*Response[Metadata])
	return response, nil
}

func GetState() (*Response[State], error) {
	resp, err := client.R().SetResult(&Response[State]{}).Get(STATE_URL)
	if err != nil {
		return nil, err
	}
	response := resp.Result().(*Response[State])
	return response, nil
}
