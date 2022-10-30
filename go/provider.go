package wheretopark

import "log"

type Provider interface {
	GetMetadata() (map[ID]Metadata, error)
	GetState() (map[ID]State, error)
}

func RunProvider(client *Client, provider Provider) error {
	metadatas, err := provider.GetMetadata()
	if err != nil {
		return err
	}
	states, err := provider.GetState()
	if err != nil {
		return err
	}

	for id, metadata := range metadatas {
		parkingLot := ParkingLot{
			Metadata: metadata,
			State:    states[id],
		}
		err := client.SetParkingLot(id, parkingLot)
		if err != nil {
			return err
		}
		log.Printf("updated parking lot of %s\n", id)
	}

	return nil
}
