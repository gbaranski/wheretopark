package wheretopark

import (
	"net/url"
	"sync"

	"github.com/rs/zerolog/log"
)

type ServerClient struct {
	baseURL *url.URL
}

func NewServerClient(baseURL *url.URL) *ServerClient {
	return &ServerClient{
		baseURL: baseURL,
	}
}

type ProviderConfig struct {
	Name string   `json:"name"`
	URL  *url.URL `json:"url"`
}

func (c *ServerClient) Providers() ([]ProviderConfig, error) {
	url := c.baseURL.JoinPath("/v1/config/providers")
	resp, err := Get[[]struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}](url, nil)
	if err != nil {
		return nil, err
	}
	providers := make([]ProviderConfig, len(*resp))
	for i, provider := range *resp {
		parsedURL, err := url.Parse(provider.Url)
		if err != nil {
			return nil, err
		}
		providers[i] = ProviderConfig{
			Name: provider.Name,
			URL:  parsedURL,
		}
	}
	return providers, nil
}

func (c *ServerClient) GetFrom(provider ProviderConfig) (map[ID]ParkingLot, error) {
	url := provider.URL.JoinPath("/parking-lots")
	parkingLots, err := Get[map[string]ParkingLot](url, nil)
	if err != nil {
		return nil, err
	}
	return *parkingLots, nil
}

func (c *ServerClient) GetFromMany(providers []ProviderConfig) (map[ID]ParkingLot, error) {
	var wg sync.WaitGroup
	var mu sync.Mutex
	allParkingLots := make(map[ID]ParkingLot)
	for _, provider := range providers {
		wg.Add(1)
		go func(provider ProviderConfig) {
			defer wg.Done()
			parkingLots, err := c.GetFrom(provider)
			if err != nil {
				log.Error().Err(err).Str("url", provider.URL.String()).Msgf("failed to fetch from %s", provider.Name)
			}
			mu.Lock()
			for id, parkingLot := range parkingLots {
				if _, ok := allParkingLots[id]; ok {
					log.Error().Str("id", id).Str("provider", provider.Name).Msg("parking lot conflict, already exists from another provider.")
					continue
				}
				allParkingLots[id] = parkingLot
			}
			mu.Unlock()
		}(provider)
	}
	wg.Wait()
	return allParkingLots, nil
}
