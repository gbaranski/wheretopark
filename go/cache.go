package wheretopark

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/allegro/bigcache/v3"
	"github.com/rs/zerolog/log"
)

type CacheProvider struct {
	c *bigcache.BigCache
}

func NewCacheProvider() (*CacheProvider, error) {
	cache, err := bigcache.New(context.Background(), bigcache.DefaultConfig(time.Minute))
	if err != nil {
		return nil, err
	}

	return &CacheProvider{
		c: cache,
	}, nil
}

func getValueFromCache[T any](cache *bigcache.BigCache, key string) (map[ID]T, bool) {
	data, err := cache.Get(key)
	if err != nil {
		if err != bigcache.ErrEntryNotFound {
			log.Error().Err(err).Msg(fmt.Sprintf("failed to get %s from cache", key))
		}
		return nil, false
	}
	log.Trace().Str("key", key).Msg(fmt.Sprintf("got `%s` from cache", data))
	var values map[ID]T
	if err := json.Unmarshal([]byte(data), &values); err != nil {
		log.Fatal().Err(err).Msg("failed to unmarshal values")
	}
	return values, true
}

func setValueToCache[T any](cache *bigcache.BigCache, key string, value map[ID]T) error {
	data, err := json.Marshal(value)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to marshal value")
	}
	log.Trace().Str("key", key).Msg(fmt.Sprintf("set `%s` to cache", data))
	return cache.Set(key, data)
}

func (p *CacheProvider) GetMetadatas(provider string) (map[ID]Metadata, bool) {
	return getValueFromCache[Metadata](p.c, fmt.Sprintf("%s/metadatas", provider))

}

func (p *CacheProvider) GetStates(provider string) (map[ID]State, bool) {
	return getValueFromCache[State](p.c, fmt.Sprintf("%s/states", provider))
}

func (p *CacheProvider) SetParkingLots(provider string, parkingLots map[ID]ParkingLot) error {
	metadatas := ExtractMetadatas(parkingLots)
	states := ExtractStates(parkingLots)
	if err := setValueToCache[Metadata](p.c, fmt.Sprintf("%s/metadatas", provider), metadatas); err != nil {
		return fmt.Errorf("failed to set metadatas")
	}
	if err := setValueToCache[State](p.c, fmt.Sprintf("%s/states", provider), states); err != nil {
		return fmt.Errorf("failed to set metadatas")
	}
	return nil
}
