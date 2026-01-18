package diyanet

import (
	"encoding/json"
	"fmt"
)

const apiURLCountries = apiURLPrefix + "api/Place/Countries"

// Country represents a country as returned by the Diyanet Awqat Salah API.
type Country struct {
	// client is the Diyanet Awqat Salah API client.
	client *Client
	// Id is the unique identifier for the country.
	Id int
	// Code is the code of the country.
	Code string
	// Name is the name of the country.
	Name string
}

// GetCountries retrieves the list of countries from the Diyanet Awqat Salah API.
func (c *Client) GetCountries() ([]Country, error) {
	resp, err := c.httpClient.Get(apiURLCountries)
	if err != nil {
		return nil, fmt.Errorf(errorPrefix+"unable to get countries: %w", err)
	}
	defer resp.Body.Close()

	var result Result[[]Country]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf(errorPrefix+"unable to decode countries response: %w", err)
	}
	if !result.Ok {
		return nil, fmt.Errorf(errorPrefix+"API error retrieving countries: %s", result.Error)
	}

	for i := range result.Data {
		result.Data[i].client = c
	}

	return result.Data, nil
}
