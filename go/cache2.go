package wheretopark

import (
	"context"
	"encoding/json"
	"fmt"
	"hash/fnv"

	"github.com/allegro/bigcache/v3"
	"github.com/rs/zerolog/log"
)

type Cache2 struct {
	i *internalCache
}

func NewCache2() (*Cache2, error) {
	cache, err := bigcache.New(context.Background(), bigcache.DefaultConfig(CacheTTL))
	if err != nil {
		return nil, err
	}

	return &Cache2{
		i: cache,
	}, nil
}

func (c *Cache2) GetParkingLotsOrUpdate(source string, or func() (map[ID]ParkingLot, error)) (map[ID]ParkingLot, error) {
	parkingLots := c.GetParkingLots(source)
	if parkingLots != nil {
		return parkingLots, nil
	}
	newParkingLots, err := or()
	if err != nil {
		return nil, fmt.Errorf("or() failed: %w", err)
	}
	c.SetParkingLots(source, newParkingLots)
	return newParkingLots, nil
}

func (c *Cache2) GetParkingLots(source string) map[ID]ParkingLot {
	data, err := c.i.Get(source)
	if err != nil {
		if err != bigcache.ErrEntryNotFound {
			log.Error().Err(err).Str("source", source).Msg("failed to get value from cache")
		}
		log.Debug().Str("key", source).Msg("cache miss")
		return nil
	}

	hash := fnv.New32a()
	hash.Write(data)

	log.Trace().Str("source", source).Uint32("sum", hash.Sum32()).Msg(fmt.Sprintf("got `%s` from cache", data))
	var value map[ID]ParkingLot
	if err := json.Unmarshal([]byte(data), &value); err != nil {
		log.Fatal().Err(err).Msg("failed to unmarshal values")
	}
	return value
}

func (c *Cache2) SetParkingLots(source string, parkingLots map[ID]ParkingLot) error {
	data, err := json.Marshal(parkingLots)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to marshal value")
	}
	hash := fnv.New32a()
	hash.Write(data)
	log.Trace().Str("source", source).Uint32("sum", hash.Sum32()).Msg(fmt.Sprintf("set `%s` to cache", data))
	return c.i.Set(source, data)
}

func (c *Cache2) UpdateParkingLots(source string, parkingLots map[ID]ParkingLot) error {
	data, err := c.i.Get(source)
	found := true
	if err != nil {
		if err == bigcache.ErrEntryNotFound {
			found = false
		} else {
			return err
		}
	}

	var value map[ID]ParkingLot
	if found {
		if err := json.Unmarshal([]byte(data), &value); err != nil {
			log.Fatal().Str("data", string(data)).Err(err).Msg("failed to unmarshal values")
		}
		for id, parkingLot := range parkingLots {
			value[id] = parkingLot
		}
	} else {
		value = parkingLots
	}

	data, err = json.Marshal(value)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to marshal value")
	}
	hash := fnv.New32a()
	hash.Write(data)
	log.Trace().Str("source", source).Uint32("sum", hash.Sum32()).Msg(fmt.Sprintf("update `%s` to cache", data))
	return c.i.Set(source, data)
}

func (c *Cache2) Close() error {
	return c.i.Close()
}
