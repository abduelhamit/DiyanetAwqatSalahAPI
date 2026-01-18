package diyanet

import (
	"encoding/json"
	"fmt"
	"time"
)

const apiURLPrayerTimeDaily = apiURLPrefix + "api/PrayerTime/Daily/%d"
const apiURLPrayerTimeWeekly = apiURLPrefix + "api/PrayerTime/Weekly/%d"
const apiURLPrayerTimeMonthly = apiURLPrefix + "api/PrayerTime/Monthly/%d"
const apiURLPrayerTimeRamadan = apiURLPrefix + "api/PrayerTime/Ramadan/%d"

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
func (c *City) GetPrayerTimeDaily(timezone *time.Location) ([]PrayerTime, error) {
	url := fmt.Sprintf(apiURLPrayerTimeDaily, c.Id)
	resp, err := c.client.httpClient.Get(url)
	if err != nil {
		return nil,
			fmt.Errorf(errorPrefix+"unable to get daily prayer time for city %s (%d – %s): %w",
				c.Name, c.Id, c.Code, err)
	}
	defer resp.Body.Close()

	var result Result[[]PrayerTime]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil,
			fmt.Errorf(errorPrefix+"unable to decode daily prayer time response for city %s (%d – %s): %w",
				c.Name, c.Id, c.Code, err)
	}
	if !result.Ok {
		return nil,
			fmt.Errorf(errorPrefix+"API error retrieving daily prayer time for city %s (%d – %s): %s",
				c.Name, c.Id, c.Code, result.Error)
	}

	for i := range result.Data {
		result.Data[i].fixGregorianDate(timezone)
	}

	return result.Data, nil
}

// GetPrayerTimeWeekly retrieves the weekly prayer times for a given city ID from the Diyanet Awqat Salah API.
// If a timezone is provided, the GregorianDate field will be adjusted to that timezone.
// If timezone is nil, the GregorianDate will be set to a fixed zone based on the GMT offset provided by the API.
func (c *City) GetPrayerTimeWeekly(cityID int, timezone *time.Location) ([]PrayerTime, error) {
	url := fmt.Sprintf(apiURLPrayerTimeWeekly, cityID)
	resp, err := c.client.httpClient.Get(url)
	if err != nil {
		return nil,
			fmt.Errorf(errorPrefix+"unable to get weekly prayer time for city %s (%d – %s): %w",
				c.Name, c.Id, c.Code, err)
	}
	defer resp.Body.Close()

	var result Result[[]PrayerTime]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil,
			fmt.Errorf(errorPrefix+"unable to decode weekly prayer time response for city %s (%d – %s): %w",
				c.Name, c.Id, c.Code, err)
	}
	if !result.Ok {
		return nil,
			fmt.Errorf(errorPrefix+"API error retrieving weekly prayer time for city %s (%d – %s): %s",
				c.Name, c.Id, c.Code, result.Error)
	}

	for i := range result.Data {
		result.Data[i].fixGregorianDate(timezone)
	}

	return result.Data, nil
}

// GetPrayerTimeMonthly retrieves the monthly prayer times for a given city ID from the Diyanet Awqat Salah API.
// If a timezone is provided, the GregorianDate field will be adjusted to that timezone.
// If timezone is nil, the GregorianDate will be set to a fixed zone based on the GMT offset provided by the API.
func (c *City) GetPrayerTimeMonthly(cityID int, timezone *time.Location) ([]PrayerTime, error) {
	url := fmt.Sprintf(apiURLPrayerTimeMonthly, cityID)
	resp, err := c.client.httpClient.Get(url)
	if err != nil {
		return nil,
			fmt.Errorf(errorPrefix+"unable to get monthly prayer time for city %s (%d – %s): %w",
				c.Name, c.Id, c.Code, err)
	}
	defer resp.Body.Close()

	var result Result[[]PrayerTime]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil,
			fmt.Errorf(errorPrefix+"unable to decode monthly prayer time response for city %s (%d – %s): %w",
				c.Name, c.Id, c.Code, err)
	}
	if !result.Ok {
		return nil,
			fmt.Errorf(errorPrefix+"API error retrieving monthly prayer time for city %s (%d – %s): %s",
				c.Name, c.Id, c.Code, result.Error)
	}

	for i := range result.Data {
		result.Data[i].fixGregorianDate(timezone)
	}

	return result.Data, nil
}

// GetPrayerTimeRamadan retrieves the Ramadan prayer times for a given city ID from the Diyanet Awqat Salah API.
// If a timezone is provided, the GregorianDate field will be adjusted to that timezone.
// If timezone is nil, the GregorianDate will be set to a fixed zone based on the GMT offset provided by the API.
func (c *City) GetPrayerTimeRamadan(cityID int, timezone *time.Location) ([]PrayerTime, error) {
	url := fmt.Sprintf(apiURLPrayerTimeRamadan, cityID)
	resp, err := c.client.httpClient.Get(url)
	if err != nil {
		return nil,
			fmt.Errorf(errorPrefix+"unable to get Ramadan prayer time for city %s (%d – %s): %w",
				c.Name, c.Id, c.Code, err)
	}
	defer resp.Body.Close()

	var result Result[[]PrayerTime]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil,
			fmt.Errorf(errorPrefix+"unable to decode Ramadan prayer time response for city %s (%d – %s): %w",
				c.Name, c.Id, c.Code, err)
	}
	if !result.Ok {
		return nil,
			fmt.Errorf(errorPrefix+"API error retrieving Ramadan prayer time for city %s (%d – %s): %s",
				c.Name, c.Id, c.Code, result.Error)
	}

	for i := range result.Data {
		result.Data[i].fixGregorianDate(timezone)
	}

	return result.Data, nil
}
