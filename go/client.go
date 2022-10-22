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

func (c *Client) GetState(id ID) (*State, error) {
	state, err := c.database.Select(fmt.Sprintf("states:%s", id))
	if err != nil {
		return nil, err
	}
	v := state.(map[string]any)
	fmt.Printf("v: %+v", v)
	availableSpotsRaw := v["available-spots"].(map[string]any)
	availableSpots := make(map[SpotType]uint)
	for key, value := range availableSpotsRaw {
		availableSpots[key] = uint(value.(float64))
	}
	return &State{
		LastUpdated:    v["last-updated"].(string),
		AvailableSpots: availableSpots,
	}, nil
	// return err
}
