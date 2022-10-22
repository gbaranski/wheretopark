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

	for id, metadata := range metadatas {
		err := client.SetMetadata(id, metadata)
		if err != nil {
			return err
		}
		log.Printf("updated metadata of %s\n", id)
	}

	states, err := provider.GetState()
	if err != nil {
		return err
	}
	for id, state := range states {
		err := client.SetState(id, state)
		if err != nil {
			return err
		}
		log.Printf("updated state of %s\n", id)
	}

	return nil
}
