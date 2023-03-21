package exporter

import (
	wheretopark "wheretopark/go"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog/log"
)

type Collector struct {
	client wheretopark.Client
}

func NewCollector(client wheretopark.Client) Collector {
	return Collector{
		client: client,
	}
}

var (
	labels             = []string{"id", "name"}
	availableSpotsDesc = prometheus.NewDesc("wheretopark_parking_lot_available_spots", "Number of available spots", labels, nil)
	totalSpotsDesc     = prometheus.NewDesc("wheretopark_parking_lot_total_spots", "Number of total spots", labels, nil)
)

func (c Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- availableSpotsDesc
	ch <- totalSpotsDesc
}

func (c Collector) Collect(ch chan<- prometheus.Metric) {
	parkingLots, err := c.client.GetAllParkingLots()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to fetch all parking lots")
	}

	for id, parkingLot := range parkingLots {
		labelValues := []string{id, parkingLot.Metadata.Name}
		ch <- prometheus.MustNewConstMetric(
			availableSpotsDesc,
			prometheus.GaugeValue,
			float64(parkingLot.State.AvailableSpots["CAR"]),
			labelValues...,
		)
		ch <- prometheus.MustNewConstMetric(
			totalSpotsDesc,
			prometheus.GaugeValue,
			float64(parkingLot.Metadata.TotalSpots["CAR"]),
			labelValues...,
		)
	}
}
