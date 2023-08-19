package glasgow

import (
	"time"

	"github.com/go-resty/resty/v2"
)

type PointCoordinates struct {
	Latitude  float64 `json:"d2lm$latitude,string"`
	Longitude float64 `json:"d2lm$longitude,string"`
}

type PointByCoordinates struct {
	PointByCoordinates PointCoordinates `json:"d2lm$pointCoordinates"`
}

type LocationContainedInGroup struct {
	PointByCoordinates PointByCoordinates `json:"d2lm$pointByCoordinates"`
}

type GroupOfLocations struct {
	LocationContainedInGroup LocationContainedInGroup `json:"d2lm$locationContainedInGroup"`
}

type SituationRecord struct {
	ID               string           `json:"@id"`
	Identity         string           `json:"d2lm$carParkIdentity"`
	DateTime         string           `json:"d2lm$situationRecordVersionTime"`
	OccupiedSpaces   int              `json:"d2lm$occupiedSpaces,string"`
	TotalCapacity    uint             `json:"d2lm$totalCapacity,string"`
	GroupOfLocations GroupOfLocations `json:"d2lm$groupOfLocations"`
}

type SituationItem struct {
	Record SituationRecord `json:"d2lm$situationRecord"`
}

type PayloadPublication struct {
	SituationItems []SituationItem `json:"d2lm$situation"`
}

type LogicalModel struct {
	PayloadPublication PayloadPublication `json:"d2lm$payloadPublication"`
}

type Response struct {
	LogicalModel LogicalModel `json:"d2lm$d2LogicalModel"`
}

const (
	// https://developer.glasgow.gov.uk/api-details#api=55c36a318b3a0306f0009483&operation=563cea91aab82f1168298575
	DATA_URL = "https://api.glasgow.gov.uk/datextraffic/carparks?format=json"
	API_KEY  = "ccaa1e24db6e4a9bb791f99433cc7ab7"
)

var client = resty.New()

func init() {
	client.GetClient().Timeout = 10 * time.Second
}

func GetData() (*Response, error) {
	resp, err := client.R().SetHeader("Ocp-Apim-Subscription-Key", API_KEY).SetHeader("Cache-Control", "no-cache").SetResult(&Response{}).Get(DATA_URL)
	if err != nil {
		return nil, err
	}
	response := resp.Result().(*Response)
	return response, nil
}
