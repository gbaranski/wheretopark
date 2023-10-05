package glasgow

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
