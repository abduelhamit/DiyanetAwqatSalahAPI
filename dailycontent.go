package diyanet

import (
	"encoding/json"
	"fmt"
)

const apiURLDailyContent = apiURLPrefix + "api/DailyContent"

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
