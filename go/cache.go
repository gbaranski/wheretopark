package wheretopark

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/allegro/bigcache/v3"
	"github.com/rs/zerolog/log"
)

func getValueFromCache[T any](cache *bigcache.BigCache, key string) *T {
	data, err := cache.Get(key)
	if err != nil {
		if err != bigcache.ErrEntryNotFound {
			log.Error().Err(err).Msg(fmt.Sprintf("failed to get %s from cache", key))
		}
		return nil
	}
	log.Trace().Str("key", key).Msg(fmt.Sprintf("got `%s` from cache", data))
	var values *T
	if err := json.Unmarshal([]byte(data), &values); err != nil {
		log.Fatal().Err(err).Msg("failed to unmarshal values")
	}
	return values
}

func setValueToCache[T any](cache *bigcache.BigCache, key string, value T) error {
	data, err := json.Marshal(value)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to marshal value")
	}
	log.Trace().Str("key", key).Msg(fmt.Sprintf("set `%s` to cache", data))
	return cache.Set(key, data)
}

type SingularCache struct {
	c *bigcache.BigCache
}

func NewSingularCache() (*SingularCache, error) {
	cache, err := bigcache.New(context.Background(), bigcache.DefaultConfig(time.Minute*5))
	if err != nil {
		return nil, err
	}

	return &SingularCache{
		c: cache,
	}, nil
}

func (p *SingularCache) GetState(id ID) *State {
	return getValueFromCache[State](p.c, fmt.Sprintf("%s/state", id))
}

func (p *SingularCache) SetState(id ID, state State) error {
	return setValueToCache[State](p.c, fmt.Sprintf("%s/state", id), state)
}

type PluralCache struct {
	c *bigcache.BigCache
}

func NewPluralCache() (*PluralCache, error) {
	cache, err := bigcache.New(context.Background(), bigcache.DefaultConfig(time.Minute*5))
	if err != nil {
		return nil, err
	}

	return &PluralCache{
		c: cache,
	}, nil
}

func (p *PluralCache) GetMetadatas(provider string) map[ID]Metadata {
	metadatas := getValueFromCache[map[ID]Metadata](p.c, fmt.Sprintf("%s/metadatas", provider))
	if metadatas != nil {
		return *metadatas
	} else {
		return nil
	}
}

func (p *PluralCache) GetStates(provider string) map[ID]State {
	states := getValueFromCache[map[ID]State](p.c, fmt.Sprintf("%s/states", provider))
	if states != nil {
		return *states
	} else {
		return nil
	}
}

func (p *PluralCache) SetParkingLots(provider string, parkingLots map[ID]ParkingLot) error {
	metadatas := ExtractMetadatas(parkingLots)
	states := ExtractStates(parkingLots)
	if err := setValueToCache[map[ID]Metadata](p.c, fmt.Sprintf("%s/metadatas", provider), metadatas); err != nil {
		return fmt.Errorf("failed to set metadatas")
	}
	if err := setValueToCache[map[ID]State](p.c, fmt.Sprintf("%s/states", provider), states); err != nil {
		return fmt.Errorf("failed to set metadatas")
	}
	return nil
}
