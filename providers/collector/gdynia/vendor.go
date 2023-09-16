package gdynia

import geojson "github.com/paulmach/go.geojson"

type Metadata struct {
	ParkingLots []struct {
		ID             int               `json:"id"`
		Code           string            `json:"code"`
		Name           string            `json:"name"`
		Address        string            `json:"address"`
		StreetEntrance string            `json:"streetEntrance"`
		Location       *geojson.Geometry `json:"location"`
		LastUpdate     string            `json:"lastUpdate"`
	} `json:"parkings"`
}

type State = []struct {
	ID         int    `json:"id"`
	ParkingID  int    `json:"parkingId"`
	Capacity   int    `json:"capacity"`
	FreePlaces int    `json:"freePlaces"`
	InsertTime string `json:"insertTime"`
}
