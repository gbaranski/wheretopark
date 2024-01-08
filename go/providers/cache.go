package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"time"
	wheretopark "wheretopark/go"

	"github.com/allegro/bigcache/v3"
	"github.com/rs/zerolog/log"
)

type internalCache = bigcache.BigCache

var CacheTTL = time.Minute * 5

type Cache struct {
	i *internalCache
}

func NewCache() (*Cache, error) {
	cache, err := bigcache.New(context.Background(), bigcache.DefaultConfig(CacheTTL))
	if err != nil {
		return nil, err
	}

	return &Cache{
		i: cache,
	}, nil
}

func (c *Cache) GetParkingLotOr(id string, or func() (wheretopark.ParkingLot, error)) (wheretopark.ParkingLot, error) {
	parkingLot := c.GetParkingLot(id)
	if parkingLot != nil {
		return *parkingLot, nil
	}
	newParkingLot, err := or()
	if err != nil {
		return wheretopark.ParkingLot{}, fmt.Errorf("or() failed: %w", err)
	}
	c.SetParkingLot(id, &newParkingLot)
	return newParkingLot, nil
}

func (c *Cache) GetParkingLot(id wheretopark.ID) *wheretopark.ParkingLot {
	data, err := c.i.Get(id)
	if err != nil {
		if err != bigcache.ErrEntryNotFound {
			log.Error().Err(err).Str("id", id).Msg("failed to get value from cache")
		}
		log.Debug().Str("key", id).Msg("cache miss")
		return nil
	}

	hash := fnv.New32a()
	hash.Write(data)

	log.Trace().Str("id", id).Uint32("sum", hash.Sum32()).Msg(fmt.Sprintf("got `%s` from cache", data))
	var value wheretopark.ParkingLot
	if err := json.Unmarshal([]byte(data), &value); err != nil {
		log.Fatal().Err(err).Msg("failed to unmarshal values")
	}
	return &value
}

func (c *Cache) GetAllParkingLots() map[wheretopark.ID]wheretopark.ParkingLot {
	if c.i.Len() == 0 {
		return nil
	}
	iterator := c.i.Iterator()
	parkingLots := make(map[wheretopark.ID]wheretopark.ParkingLot, c.i.Len())
	for iterator.SetNext() {
		current, err := iterator.Value()
		if err != nil {
			log.Error().Err(err).Msg("failed to iterate over cache")
			continue
		}

		data := current.Value()
		id := wheretopark.ID(current.Key())
		hash := fnv.New32a()
		hash.Write(data)

		log.Trace().Str("id", id).Uint32("sum", hash.Sum32()).Msg(fmt.Sprintf("got `%s` from cache", data))

		var value wheretopark.ParkingLot
		if err := json.Unmarshal([]byte(data), &value); err != nil {
			log.Fatal().Err(err).Msg("failed to unmarshal values")
		}
		parkingLots[id] = value
	}

	return parkingLots
}

func (c *Cache) SetParkingLot(id wheretopark.ID, parkingLot *wheretopark.ParkingLot) error {
	data, err := json.Marshal(parkingLot)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to marshal value")
	}
	hash := fnv.New32a()
	hash.Write(data)
	log.Trace().Str("id", id).Uint32("sum", hash.Sum32()).Msg(fmt.Sprintf("set `%s` to cache", data))
	return c.i.Set(id, data)
}

func (c *Cache) Close() error {
	return c.i.Close()
}
