package cctv

import (
	wheretopark "wheretopark/go"
)

type Provider struct {
	configuration Configuration
	model         Model
	saver         Saver
}

func NewProvider(model Model, saver Saver) (*Provider, error) {
	return &Provider{
		configuration: DefaultConfiguration,
		model:         model,
		saver:         saver,
	}, nil
}

func (p *Provider) Sources() map[string]wheretopark.Source {
	sources := make(map[string]wheretopark.Source, len(p.configuration.ParkingLots))
	for id, parkingLot := range p.configuration.ParkingLots {
		source := Source{
			id:       id,
			metadata: parkingLot.Metadata,
			cameras:  parkingLot.Cameras,
			model:    p.model,
			saver:    p.saver,
		}
		sources[id] = wheretopark.NewSequentialSourceProxy(&source)
	}
	return sources
}

func (p *Provider) Close() {
	p.model.Close()
}
