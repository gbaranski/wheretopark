package cctv

import (
	wheretopark "wheretopark/go"
)

type Provider struct {
	configuration Configuration
	model         *Model
	saver         Saver
}

func NewProvider(modelPath string, saveBasePath *string, saveItems []SaveItem, saveIDs []wheretopark.ID) (*Provider, error) {
	model := NewModel(modelPath)
	saver := NewSaver(saveBasePath, saveItems, saveIDs)
	return &Provider{
		configuration: DefaultConfiguration,
		model:         model,
		saver:         saver,
	}, nil
}

func (p *Provider) Sources() map[string]wheretopark.Source {
	sources := make(map[string]wheretopark.Source, len(p.configuration.ParkingLots))
	for _, parkingLot := range p.configuration.ParkingLots {
		id := wheretopark.GeometryToID(parkingLot.Geometry)
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
