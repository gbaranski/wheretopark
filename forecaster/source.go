package forecaster

import "time"

type Source interface {
	Name() string
	Load() (map[string]*ParkingMeter, error)
}

type ParkingMeter struct {
	Name          string
	OccupancyData map[time.Time]int
}

const Interval time.Duration = time.Minute * 15

func (p *ParkingMeter) AddOccupancy(start time.Time, end time.Time) {
	currentInterval := start.Truncate(Interval)
	endInterval := end.Truncate(Interval)

	for currentInterval.Before(endInterval) {
		p.OccupancyData[currentInterval]++
		currentInterval = currentInterval.Add(Interval)
	}
}

func (p *ParkingMeter) TotalSpots() uint {
	count := uint(0)
	for _, value := range p.OccupancyData {
		count = max(count, uint(value))
	}
	return count
}
