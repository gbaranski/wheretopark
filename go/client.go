package wheretopark

import (
	"fmt"
	"net/url"

	"github.com/goccy/go-json"
	"github.com/surrealdb/surrealdb.go"
)

type Client struct {
	database *surrealdb.DB
}

func NewClient(url *url.URL, namespace, databaseName string) (*Client, error) {
	endpoint := url.JoinPath("rpc")
	database, err := surrealdb.New(endpoint.String())
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

func (c *Client) getAll(table string) ([]any, error) {
	entries, err := c.database.Select(table)
	if err != nil {
		return nil, err
	}
	v := entries.([]any)
	return v, nil
}

func (c *Client) delete(thing string) error {
	_, err := c.database.Delete(thing)
	return err
}

func parkingLotReference(id ID) string {
	return fmt.Sprintf("parking_lot:%s", id)
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

func (c *Client) SetParkingLot(id ID, parkingLot ParkingLot) error {
	data, err := toMap(parkingLot)
	if err != nil {
		return err
	}
	return c.set(parkingLotReference(id), data)
}

func (c *Client) ExistsParkingLot(id ID) (bool, error) {
	return c.exists(parkingLotReference(id))
}

func (c *Client) GetParkingLot(id ID) (*ParkingLot, error) {
	v, err := c.get(parkingLotReference(id))
	if err != nil {
		return nil, err
	}
	if v == nil {
		return nil, nil
	}
	parkingLot, err := mapTo[ParkingLot](v)
	if err != nil {
		return nil, err
	}
	return parkingLot, nil
}

func (c *Client) DeleteParkingLot(id ID) error {
	return c.delete(parkingLotReference(id))
}

func (c *Client) GetAllParkingLots() (map[ID]ParkingLot, error) {
	v, err := c.getAll("parking_lot")
	if err != nil {
		return nil, err
	}
	parkingLots := make(map[ID]ParkingLot, len(v))
	for _, data := range v {
		data := data.(map[string]any)
		id := data["id"].(string)
		id = id[len("parking_lot:"):]
		parkingLot, err := mapTo[ParkingLot](data)
		if err != nil {
			return nil, err
		}
		parkingLots[id] = *parkingLot
	}
	return parkingLots, nil
}

func (c *Client) Close() {
	c.database.Close()
}
