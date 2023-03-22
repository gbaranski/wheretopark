package gdynia

import (
	"github.com/go-resty/resty/v2"
	geojson "github.com/paulmach/go.geojson"
)

type Metadata struct {
	ID             int               `json:"id"`
	Code           string            `json:"code"`
	Name           string            `json:"name"`
	Address        string            `json:"address"`
	StreetEntrance string            `json:"streetEntrance"`
	Location       *geojson.Geometry `json:"location"`
	LastUpdate     string            `json:"lastUpdate"`
}

type MetadataResponse struct {
	Parkings []Metadata `json:"parkings"`
}

type State struct {
	ID         int    `json:"id"`
	ParkingID  int    `json:"parkingId"`
	Capacity   int    `json:"capacity"`
	FreePlaces int    `json:"freePlaces"`
	InsertTime string `json:"insertTime"`
}

type StateResponse = []State

const (
	METADATA_URL = "http://api.zdiz.gdynia.pl/ri/rest/parkings"
	STATE_URL    = "http://api.zdiz.gdynia.pl/ri/rest/parking_places"
)

var client = resty.New()

func GetMetadata() (*MetadataResponse, error) {
	resp, err := client.R().SetResult(&MetadataResponse{}).Get(METADATA_URL)
	if err != nil {
		return nil, err
	}
	response := resp.Result().(*MetadataResponse)
	return response, nil
}

func GetState() (*StateResponse, error) {
	resp, err := client.R().SetResult(&StateResponse{}).Get(STATE_URL)
	if err != nil {
		return nil, err
	}
	response := resp.Result().(*StateResponse)
	return response, nil
}
