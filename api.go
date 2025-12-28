package diyanet

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const apiURLDailyContent = apiURLPrefix + "api/DailyContent"
const apiURLCountries = apiURLPrefix + "api/Place/Countries"
const apiURLStates = apiURLPrefix + "api/Place/States"
const apiURLStatesByCountry = apiURLPrefix + "api/Place/States/%d"
const apiURLCities = apiURLPrefix + "api/Place/Cities"
const apiURLCitiesByState = apiURLPrefix + "api/Place/Cities/%d"
const apiURLCityDetail = apiURLPrefix + "api/Place/CityDetail/%d"
const apiURLPrayerTimeDaily = apiURLPrefix + "api/PrayerTime/Daily/%d"
const apiURLPrayerTimeWeekly = apiURLPrefix + "api/PrayerTime/Weekly/%d"
const apiURLPrayerTimeMonthly = apiURLPrefix + "api/PrayerTime/Monthly/%d"

// Client is a Diyanet Awqat Salah API client.
type Client struct {
	// httpClient is the HTTP client used to make requests.
	httpClient *http.Client
}

// NewClient creates a new Diyanet Awqat Salah API client using the provided configuration.
func (c *Config) NewClient(ctx context.Context) *Client {
	return &Client{
		httpClient: c.HTTPClient(ctx),
	}
}

// DailyContent describes a single day's devotional content—verse (ayah),
// hadith, and prayer (du'a)—along with their source references and metadata.
type DailyContent struct {
	// Id is the unique identifier for this content record.
	Id int
	// DayOfYear is the 1–366 ordinal day in the calendar year this content applies to.
	DayOfYear int
	// Verse is the textual content of the selected verse.
	Verse string
	// VerseSource identifies the verse's reference (e.g., "(Şu'arâ, 42/29)").
	VerseSource string
	// Hadith is the textual content of the selected hadith.
	Hadith string
	// HadithSource identifies the hadith's reference (e.g., "(Tirmizî, “Birr ”, 15)").
	HadithSource string
	// Pray is the textual content of the selected prayer.
	Pray string
	// PraySource identifies the prayer's source or attribution.
	PraySource string
}

// GetDailyContent retrieves the daily content from the Diyanet Awqat Salah API.
func (c *Client) GetDailyContent() (DailyContent, error) {
	resp, err := c.httpClient.Get(apiURLDailyContent)
	if err != nil {
		return DailyContent{}, fmt.Errorf(errorPrefix+"unable to get daily content: %w", err)
	}
	defer resp.Body.Close()

	var result Result[DailyContent]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return DailyContent{}, fmt.Errorf(errorPrefix+"unable to decode daily content response: %w", err)
	}
	if !result.Ok {
		return DailyContent{}, fmt.Errorf(errorPrefix+"API error retrieving daily content: %s", result.Error)
	}

	return result.Data, nil
}

// Country represents a country as returned by the Diyanet Awqat Salah API.
type Country struct {
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

	return result.Data, nil
}

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
func (c *Client) GetCityDetail(cityID int) (CityDetail, error) {
	url := fmt.Sprintf(apiURLCityDetail, cityID)
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return CityDetail{}, fmt.Errorf(errorPrefix+"unable to get city detail for city ID %d: %w", cityID, err)
	}
	defer resp.Body.Close()

	var result Result[CityDetail]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return CityDetail{}, fmt.Errorf(errorPrefix+"unable to decode city detail response for city ID %d: %w", cityID, err)
	}
	if !result.Ok {
		return CityDetail{}, fmt.Errorf(errorPrefix+"API error retrieving city detail for city ID %d: %s", cityID, result.Error)
	}

	return result.Data, nil
}

// PrayerTime represents the prayer times and related information for a specific day in a city.
type PrayerTime struct {
	// ShapeMoonURL is the URL of the moon phase image.
	ShapeMoonURL string
	// Fajr is the time for the Fajr prayer.
	Fajr string
	// Sunrise is the time for sunrise.
	Sunrise string
	// Dhuhr is the time for the Dhuhr prayer.
	Dhuhr string
	// Asr is the time for the Asr prayer.
	Asr string
	// Maghrib is the time for the Maghrib prayer.
	Maghrib string
	// Isha is the time for the Isha prayer.
	Isha string
	// AstronomicalSunset is the time for astronomical sunset.
	AstronomicalSunset string
	// AstronomicalSunrise is the time for astronomical sunrise.
	AstronomicalSunrise string
	// HijriDateShort is the short format of the Hijri date.
	HijriDateShort string
	// HijriDateLong is the long format of the Hijri date.
	HijriDateLong string
	// HijriDate is the Hijri date as a time.Time object.
	HijriDate time.Time `json:"hijriDateLongIso8601"`
	// QiblaTime is the time for Qibla.
	QiblaTime string
	// GregorianDateShort is the short format of the Gregorian date.
	GregorianDateShort string
	// GregorianDateLong is the long format of the Gregorian date.
	GregorianDateLong string
	// GregorianDate is the Gregorian date as a time.Time object.
	GregorianDate time.Time `json:"gregorianDateLongIso8601"`
	// GreenwichMeanTimeZone is the GMT offset for the location.
	GreenwichMeanTimeZone float32
}

func (pt *PrayerTime) fixGregorianDate(timezone *time.Location) {
	if timezone == nil {
		timezone = time.FixedZone(fmt.Sprintf("GMT%.2f", pt.GreenwichMeanTimeZone), int(pt.GreenwichMeanTimeZone*3600))
	}

	pt.GregorianDate = time.Date(
		pt.GregorianDate.Year(),
		pt.GregorianDate.Month(),
		pt.GregorianDate.Day(),
		0, 0, 0, 0,
		timezone,
	)
}

// GetPrayerTimeDaily retrieves the daily prayer times for a given city ID from the Diyanet Awqat Salah API.
// If a timezone is provided, the GregorianDate field will be adjusted to that timezone.
// If timezone is nil, the GregorianDate will be set to a fixed zone based on the GMT offset provided by the API.
func (c *Client) GetPrayerTimeDaily(cityID int, timezone *time.Location) ([]PrayerTime, error) {
	url := fmt.Sprintf(apiURLPrayerTimeDaily, cityID)
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf(errorPrefix+"unable to get daily prayer time for city ID %d: %w", cityID, err)
	}
	defer resp.Body.Close()

	var result Result[[]PrayerTime]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf(errorPrefix+"unable to decode daily prayer time response for city ID %d: %w", cityID, err)
	}
	if !result.Ok {
		return nil, fmt.Errorf(errorPrefix+"API error retrieving daily prayer time for city ID %d: %s", cityID, result.Error)
	}

	for i := range result.Data {
		result.Data[i].fixGregorianDate(timezone)
	}

	return result.Data, nil
}

// GetPrayerTimeWeekly retrieves the weekly prayer times for a given city ID from the Diyanet Awqat Salah API.
// If a timezone is provided, the GregorianDate field will be adjusted to that timezone.
// If timezone is nil, the GregorianDate will be set to a fixed zone based on the GMT offset provided by the API.
func (c *Client) GetPrayerTimeWeekly(cityID int, timezone *time.Location) ([]PrayerTime, error) {
	url := fmt.Sprintf(apiURLPrayerTimeWeekly, cityID)
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf(errorPrefix+"unable to get weekly prayer time for city ID %d: %w", cityID, err)
	}
	defer resp.Body.Close()

	var result Result[[]PrayerTime]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf(errorPrefix+"unable to decode weekly prayer time response for city ID %d: %w", cityID, err)
	}
	if !result.Ok {
		return nil, fmt.Errorf(errorPrefix+"API error retrieving weekly prayer time for city ID %d: %s", cityID, result.Error)
	}

	for i := range result.Data {
		result.Data[i].fixGregorianDate(timezone)
	}

	return result.Data, nil
}

// GetPrayerTimeMonthly retrieves the monthly prayer times for a given city ID from the Diyanet Awqat Salah API.
// If a timezone is provided, the GregorianDate field will be adjusted to that timezone.
// If timezone is nil, the GregorianDate will be set to a fixed zone based on the GMT offset provided by the API.
func (c *Client) GetPrayerTimeMonthly(cityID int, timezone *time.Location) ([]PrayerTime, error) {
	url := fmt.Sprintf(apiURLPrayerTimeMonthly, cityID)
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf(errorPrefix+"unable to get monthly prayer time for city ID %d: %w", cityID, err)
	}
	defer resp.Body.Close()

	var result Result[[]PrayerTime]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf(errorPrefix+"unable to decode monthly prayer time response for city ID %d: %w", cityID, err)
	}
	if !result.Ok {
		return nil, fmt.Errorf(errorPrefix+"API error retrieving monthly prayer time for city ID %d: %s", cityID, result.Error)
	}

	for i := range result.Data {
		result.Data[i].fixGregorianDate(timezone)
	}

	return result.Data, nil
}
