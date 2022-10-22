package wheretopark

import (
	"fmt"

	"github.com/surrealdb/surrealdb.go"
)

type Client struct {
	database *surrealdb.DB
}

func NewClient(url, namespace, database string) (*Client, error) {
	surrealdb, err := surrealdb.New(url)
	if err != nil {
		return nil, err
	}
	_, err = surrealdb.Use(namespace, database)
	if err != nil {
		return nil, err
	}

	return &Client{
		database: surrealdb,
	}, nil
}

func (c *Client) SignInWithPassword(user, password string) error {
	_, err := c.database.Signin(map[string]any{
		"user": user,
		"pass": password,
	})
	return err
}

func (c *Client) AddState(id ID, state State) error {
	_, err := c.database.Create(fmt.Sprintf("states:%s", id), map[string]any{
		"last-updated":    state.LastUpdated,
		"available-spots": state.AvailableSpots,
	})
	return err
}

func (c *Client) StateExists(id ID) (bool, error) {
	raw, err := c.database.Query("SELECT true as exists FROM $what", map[string]any{"what": fmt.Sprintf("states:%s", id)})
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

func (c *Client) GetState(id ID) (*State, error) {
	exists, err := c.StateExists(id)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, nil
	}
	state, err := c.database.Select(fmt.Sprintf("states:%s", id))
	if err != nil {
		return nil, err
	}
	v := state.(map[string]any)
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
	_, err := c.database.Delete(fmt.Sprintf("states:%s", id))
	return err
}
