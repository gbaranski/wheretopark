package warsaw

type FreePlacesTotal struct {
	Disabled int `json:"disabled"`
	Public   int `json:"public"`
	Electric int `json:"electric"`
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
