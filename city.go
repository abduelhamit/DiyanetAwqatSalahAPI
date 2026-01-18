package diyanet

import (
	"encoding/json"
	"fmt"
)

const apiURLCities = apiURLPrefix + "api/Place/Cities"
const apiURLCitiesByState = apiURLPrefix + "api/Place/Cities/%d"

// City represents a city as returned by the Diyanet Awqat Salah API.
type City struct {
	// Id is the unique identifier for the city.
	Id int
	// Code is the code of the city.
	Code string
	// Name is the name of the city.
	Name string
}

// GetCities retrieves the list of cities from the Diyanet Awqat Salah API.
func (c *Client) GetCities() ([]City, error) {
	resp, err := c.httpClient.Get(apiURLCities)
	if err != nil {
		return nil, fmt.Errorf(errorPrefix+"unable to get cities: %w", err)
	}
	defer resp.Body.Close()

	var result Result[[]City]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf(errorPrefix+"unable to decode cities response: %w", err)
	}
	if !result.Ok {
		return nil, fmt.Errorf(errorPrefix+"API error retrieving cities: %s", result.Error)
	}

	return result.Data, nil
}

// GetCitiesByState retrieves the list of cities for a given state ID from the Diyanet Awqat Salah API.
func (c *Client) GetCitiesByState(stateID int) ([]City, error) {
	url := fmt.Sprintf(apiURLCitiesByState, stateID)
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf(errorPrefix+"unable to get cities for state ID %d: %w", stateID, err)
	}
	defer resp.Body.Close()

	var result Result[[]City]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf(errorPrefix+"unable to decode cities response for state ID %d: %w", stateID, err)
	}
	if !result.Ok {
		return nil, fmt.Errorf(errorPrefix+"API error retrieving cities for state ID %d: %s", stateID, result.Error)
	}

	return result.Data, nil
}
