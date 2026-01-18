package diyanet

import (
	"encoding/json"
	"fmt"
)

const apiURLCityDetail = apiURLPrefix + "api/Place/CityDetail/%d"

// CityDetail represents detailed information about a city as returned by the Diyanet Awqat Salah API.
type CityDetail struct {
	// Id is the unique identifier for the city.
	Id string
	// Name is the name of the city.
	Name string
	// Code is the code of the city.
	Code string
	// GeographicQiblaAngle is the geographic Qibla angle for the city.
	GeographicQiblaAngle string
	// DistanceToKaaba is the distance to the Kaaba from the city.
	DistanceToKaaba string
	// QiblaAngle is the Qibla angle for the city.
	QiblaAngle string
	// City is the name of the city.
	City string
	// CityEn is the English name of the city.
	CityEn string
	// Country is the name of the country.
	Country string
	// CountryEn is the English name of the country.
	CountryEn string
}

// GetCityDetail retrieves detailed information about a city by its ID from the Diyanet Awqat Salah API.
func (c *City) GetCityDetail() (CityDetail, error) {
	url := fmt.Sprintf(apiURLCityDetail, c.Id)
	resp, err := c.client.httpClient.Get(url)
	if err != nil {
		return CityDetail{},
			fmt.Errorf(errorPrefix+"unable to get city detail for city %s (%d – %s): %w",
				c.Name, c.Id, c.Code, err)
	}
	defer resp.Body.Close()

	var result Result[CityDetail]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return CityDetail{},
			fmt.Errorf(errorPrefix+"unable to decode city detail response for city %s (%d – %s): %w",
				c.Name, c.Id, c.Code, err)
	}
	if !result.Ok {
		return CityDetail{},
			fmt.Errorf(errorPrefix+"API error retrieving city detail for city %s (%d – %s): %s",
				c.Name, c.Id, c.Code, result.Error)
	}

	return result.Data, nil
}
