package lacity

import (
	"sort"
	"time"
	wheretopark "wheretopark/go"
	"wheretopark/go/provider/sequential"

	"github.com/rs/zerolog/log"
	"golang.org/x/exp/slices"
	"golang.org/x/text/currency"

	geojson "github.com/paulmach/go.geojson"
)

type Provider struct {
	// mapping wheretopark.ID -> []SpaceID
	mappings map[wheretopark.ID][]string
}

func (p Provider) Config() sequential.Config {
	return sequential.DEFAULT_CONFIG
}

func (p Provider) Name() string {
	return "lacity"
}

func getGroups(vSpaces []SpaceMetadata) map[string][]SpaceMetadata {
	groups := map[string][]SpaceMetadata{}
	for _, space := range vSpaces {
		if groups[space.BlockFace] == nil {
			groups[space.BlockFace] = []SpaceMetadata{}
		}
		groups[space.BlockFace] = append(groups[space.BlockFace], space)
	}
	return groups
}

func sortGroups(groups map[string][]SpaceMetadata) {
	for _, group := range groups {
		sort.SliceStable(group, func(i, j int) bool {
			return group[i].SpaceID < group[j].SpaceID
		})
	}
}

func spaceIDs(spaces []SpaceMetadata) []string {
	ids := make([]string, len(spaces))
	for i, space := range spaces {
		ids[i] = space.SpaceID
	}
	return ids
}

var timezone *time.Location = wheretopark.MustLoadLocation("America/Los_Angeles")

func (p Provider) GetMetadatas() (map[wheretopark.ID]wheretopark.Metadata, error) {
	vSpaces, err := GetSpaceMetadatas()
	if err != nil {
		return nil, err
	}
	groups := getGroups(vSpaces)
	sortGroups(groups)

	metadatas := make(map[wheretopark.ID]wheretopark.Metadata, len(groups))
	for blockFace, spaces := range groups {
		mainSpace := spaces[0]
		rateType := mainSpace.RateType
		for _, v := range spaces {
			if v.RateType != rateType {
				log.Warn().Err(err).Str("blockFace", blockFace).Str("spaceID", v.SpaceID).Msg("not all rate types are the same")
			}
		}

		rules, err := mainSpace.Rules()
		if err != nil {
			log.Warn().Err(err).Str("blockFace", blockFace).Str("spaceID", mainSpace.SpaceID).Msg("failed to parse as rules")
			continue
		}

		metadata := wheretopark.Metadata{
			LastUpdated: nil,
			Name:        blockFace,
			Address:     blockFace,
			Geometry:    geojson.NewPointGeometry([]float64{spaces[0].Coordinate.Longitude, spaces[0].Coordinate.Latitude}),
			Resources: []string{
				"https://ladotparking.org",
			},
			TotalSpots: map[string]uint{
				wheretopark.SpotTypeCar: uint(len(spaces)),
			},
			MaxDimensions: nil,
			Features: []wheretopark.Feature{
				wheretopark.FeatureUncovered,
			},
			PaymentMethods: []wheretopark.PaymentMethod{
				wheretopark.PaymentMethodCash, wheretopark.PaymentMethodCard,
			},
			Comment: map[string]string{
				"en": "Source of data: http://www.laexpresspark.org/",
			},
			Currency: currency.USD,
			Timezone: timezone,
			Rules:    rules,
		}
		id := wheretopark.GeometryToID(metadata.Geometry)
		metadatas[id] = metadata
		p.mappings[id] = spaceIDs(spaces)
	}
	return metadatas, nil
}

func (p Provider) GetStates() (map[wheretopark.ID]wheretopark.State, error) {
	spaces, err := GetSpaceStates()
	if err != nil {
		return nil, err
	}

	states := make(map[wheretopark.ID]wheretopark.State)
	for _, space := range spaces {
		var id wheretopark.ID
		for i, spaceIDs := range p.mappings {
			if slices.Contains(spaceIDs, space.SpaceID) {
				id = i
				break
			}
		}
		if id == "" {
			continue
		}
		spaceLastUpdated, err := time.ParseInLocation("2006-01-02T15:04:05.000", space.EventTime, time.UTC)
		if err != nil {
			log.Warn().Err(err).Str("spaceID", space.SpaceID).Msg("failed to parse last update time")
			continue
		}

		var newAvailableSpots uint
		if space.OccupancyState == OccupancyStateVacant {
			newAvailableSpots = 1
		}

		if state, exists := states[id]; exists {
			states[id].AvailableSpots[wheretopark.SpotTypeCar] += newAvailableSpots
			if state.LastUpdated.Before(spaceLastUpdated) {
				state.LastUpdated = spaceLastUpdated
			}
		} else {
			states[id] = wheretopark.State{
				LastUpdated: spaceLastUpdated,
				AvailableSpots: map[string]uint{
					wheretopark.SpotTypeCar: newAvailableSpots,
				},
			}
		}
	}
	return states, nil
}

func NewProvider() (sequential.Provider, error) {
	return Provider{
		mappings: make(map[wheretopark.ID][]string),
	}, nil
}
