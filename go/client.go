package wheretopark

import (
	"fmt"

	"github.com/goccy/go-json"
	"github.com/surrealdb/surrealdb.go"
)

type Client struct {
	database *surrealdb.DB
}

func NewClient(url, namespace, databaseName string) (*Client, error) {
	database, err := surrealdb.New(url)
	if err != nil {
		return nil, err
	}
	_, err = database.Use(namespace, databaseName)
	if err != nil {
		return nil, err
	}

	return &Client{
		database: database,
	}, nil
}

func (c *Client) SignInWithPassword(user, password string) error {
	_, err := c.database.Signin(map[string]any{
		"user": user,
		"pass": password,
	})
	return err
}

func (c *Client) set(thing string, data map[string]any) error {
	_, err := c.database.Update(thing, data)
	return err
}

func (c *Client) exists(thing string) (bool, error) {
	raw, err := c.database.Query("SELECT true as exists FROM $what", map[string]any{"what": thing})
	if err != nil {
		return false, err
	}
	responses := raw.([]any)
	response := responses[0].(map[string]any)
	results := response["result"].([]any)
	if len(results) == 0 {
		return false, nil
	}
	result := results[0].(map[string]any)
	exists := result["exists"].(bool)
	return exists, nil
}

func (c *Client) get(thing string) (map[string]any, error) {
	exists, err := c.exists(thing)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, nil
	}
	state, err := c.database.Select(thing)
	if err != nil {
		return nil, err
	}
	v := state.(map[string]any)
	return v, nil
}

func (c *Client) delete(thing string) error {
	_, err := c.database.Delete(thing)
	return err
}

func metadataReference(id ID) string {
	return fmt.Sprintf("metadatas:%s", id)
}

func toMap(v any) (map[string]any, error) {
	bytes, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	var value map[string]any
	err = json.Unmarshal(bytes, &value)
	if err != nil {
		return nil, err
	}
	return value, nil
}

func mapTo[T any](v map[string]any) (*T, error) {
	bytes, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	var value T
	err = json.Unmarshal(bytes, &value)
	if err != nil {
		return nil, err
	}
	return &value, nil
}

func (c *Client) SetMetadata(id ID, metadata Metadata) error {
	data, err := toMap(metadata)
	if err != nil {
		return err
	}
	return c.set(metadataReference(id), data)
}

func (c *Client) ExistsMetadata(id ID) (bool, error) {
	return c.exists(metadataReference(id))
}

func (c *Client) GetMetadata(id ID) (*Metadata, error) {
	v, err := c.get(metadataReference(id))
	if err != nil {
		return nil, err
	}
	if v == nil {
		return nil, nil
	}
	metadata, err := mapTo[Metadata](v)
	if err != nil {
		return nil, err
	}
	return metadata, nil
}

func (c *Client) DeleteMetadata(id ID) error {
	return c.delete(metadataReference(id))
}

func stateReference(id ID) string {
	return fmt.Sprintf("states:%s", id)
}

func (c *Client) SetState(id ID, state State) error {
	return c.set(stateReference(id), map[string]any{
		"last-updated":    state.LastUpdated,
		"available-spots": state.AvailableSpots,
	})
}

func (c *Client) ExistsState(id ID) (bool, error) {
	return c.exists(stateReference(id))
}

func (c *Client) GetState(id ID) (*State, error) {
	v, err := c.get(stateReference(id))
	if err != nil {
		return nil, err
	}
	if v == nil {
		return nil, nil
	}
	availableSpotsRaw := v["available-spots"].(map[string]any)
	availableSpots := make(map[SpotType]uint)
	for key, value := range availableSpotsRaw {
		availableSpots[key] = uint(value.(float64))
	}
	return &State{
		LastUpdated:    v["last-updated"].(string),
		AvailableSpots: availableSpots,
	}, nil
}

func (c *Client) DeleteState(id ID) error {
	return c.delete(stateReference(id))
}
