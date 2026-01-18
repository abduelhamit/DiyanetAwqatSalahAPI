package diyanet

import (
	"encoding/json"
	"fmt"
)

const apiURLStates = apiURLPrefix + "api/Place/States"
const apiURLStatesByCountry = apiURLPrefix + "api/Place/States/%d"

// State represents a state or province as returned by the Diyanet Awqat Salah API.
type State struct {
	// Id is the unique identifier for the state.
	Id int
	// Code is the code of the state.
	Code string
	// Name is the name of the state.
	Name string
}

// GetStates retrieves the list of states from the Diyanet Awqat Salah API.
func (c *Client) GetStates() ([]State, error) {
	resp, err := c.httpClient.Get(apiURLStates)
	if err != nil {
		return nil, fmt.Errorf(errorPrefix+"unable to get states: %w", err)
	}
	defer resp.Body.Close()

	var result Result[[]State]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf(errorPrefix+"unable to decode states response: %w", err)
	}
	if !result.Ok {
		return nil, fmt.Errorf(errorPrefix+"API error retrieving states: %s", result.Error)
	}

	return result.Data, nil
}

// GetStatesByCountry retrieves the list of states for a given country ID from the Diyanet Awqat Salah API.
func (c *Client) GetStatesByCountry(countryID int) ([]State, error) {
	url := fmt.Sprintf(apiURLStatesByCountry, countryID)
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf(errorPrefix+"unable to get states for country ID %d: %w", countryID, err)
	}
	defer resp.Body.Close()

	var result Result[[]State]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf(errorPrefix+"unable to decode states response for country ID %d: %w", countryID, err)
	}
	if !result.Ok {
		return nil, fmt.Errorf(errorPrefix+"API error retrieving states for country ID %d: %s", countryID, result.Error)
	}

	return result.Data, nil
}
