package warsaw

import (
	"time"

	"github.com/go-resty/resty/v2"
)

type FreePlacesTotal struct {
	Disabled uint `json:"disabled"`
	Public   uint `json:"public"`
	Electric uint `json:"electric"`
}

type TotalPlaces struct {
	Disabled uint `json:"disabled"`
	Standard uint `json:"standard"`
	Electric uint `json:"electric"`
}

type Dimension struct {
	Width  string `json:"Width"`
	Length string `json:"Length"`
}

type Park struct {
	Name            string          `json:"name"`
	FreePlacesTotal FreePlacesTotal `json:"free_places_total"`
	TotalPlaces     []TotalPlaces   `json:"total_places"`
	Longitude       float64         `json:"longitude,string"`
	Latitude        float64         `json:"latitude,string"`
	Address         string          `json:"adress"`
	Dimensions      []Dimension     `json:"dimensions"`
}

type Result struct {
	Status int `json:"Status"`
	// ISO 8601
	Timestamp string `json:"Timestamp"`
	Parks     []Park `json:"Parks"`
}

type Response struct {
	Result Result `json:"result"`
}

const (
	// API key from https://api.um.warszawa.pl
	API_KEY  = "8840e9a0-0a23-4d9a-90f8-8c9a49b88e3b"
	DATA_URL = "https://api.um.warszawa.pl/api/action/parking_get_list/?apikey=" + API_KEY
)

var client = resty.New()

func init() {
	client.GetClient().Timeout = 10 * time.Second
}

func GetData() (*Response, error) {
	resp, err := client.R().SetResult(&Response{}).Get(DATA_URL)
	if err != nil {
		return nil, err
	}
	response := resp.Result().(*Response)
	return response, nil
}
