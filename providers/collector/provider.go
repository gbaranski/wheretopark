package collector

import (
	wheretopark "wheretopark/go"
	"wheretopark/providers/collector/gdansk"
	"wheretopark/providers/collector/gdynia"
	"wheretopark/providers/collector/glasgow"
	"wheretopark/providers/collector/lacity"
	"wheretopark/providers/collector/poznan"
	"wheretopark/providers/collector/warsaw"
)

type Provider struct {
	sources map[string]wheretopark.Source
}

func NewProvider() (*Provider, error) {
	sources := map[string]wheretopark.Source{
		"gdansk":  wheretopark.NewSequentialSourceProxy(gdansk.New()),
		"gdynia":  wheretopark.NewSequentialSourceProxy(gdynia.New()),
		"glasgow": glasgow.New(),
		"lacity":  wheretopark.NewSequentialSourceProxy(lacity.New()),
		"poznan":  poznan.New(),
		"warsaw":  warsaw.New(),
	}

	return &Provider{
		sources: sources,
	}, nil
}

func (p *Provider) Sources() map[string]wheretopark.Source {
	return p.sources
}
