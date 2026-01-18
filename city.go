package diyanet

import (
	"encoding/json"
	"fmt"
)

const apiURLCities = apiURLPrefix + "api/Place/Cities"
const apiURLCitiesByState = apiURLPrefix + "api/Place/Cities/%d"

// City represents a city as returned by the Diyanet Awqat Salah API.
type City struct {
	// client is the Diyanet Awqat Salah API client.
	client *Client
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

	for i := range result.Data {
		result.Data[i].client = c
	}

	return result.Data, nil
}

// GetCities retrieves the list of cities for a given state from the Diyanet Awqat Salah API.
func (s *State) GetCities() ([]City, error) {
	url := fmt.Sprintf(apiURLCitiesByState, s.Id)
	resp, err := s.client.httpClient.Get(url)
	if err != nil {
		return nil,
			fmt.Errorf(errorPrefix+"unable to get cities for state %s (%d – %s): %w",
				s.Name, s.Id, s.Code, err)
	}
	defer resp.Body.Close()

	var result Result[[]City]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil,
			fmt.Errorf(errorPrefix+"unable to decode cities response for state %s (%d – %s): %w",
				s.Name, s.Id, s.Code, err)
	}
	if !result.Ok {
		return nil,
			fmt.Errorf(errorPrefix+"API error retrieving cities for state %s (%d – %s): %s",
				s.Name, s.Id, s.Code, result.Error)
	}

	for i := range result.Data {
		result.Data[i].client = s.client
	}

	return result.Data, nil
}
