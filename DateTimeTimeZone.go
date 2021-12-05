package msgraph

import (
	"encoding/json"
	"fmt"
	"time"
)

// DateTimeTimeZone represents a proposed time for CalendarEvent
type DateTimeTimeZone struct {
	DateTime  *time.Time  `json:"dateTime,omitempty"`
	TimeZone  *string   `json:"timeZone,omitempty"`
}

func (s DateTimeTimeZone) String() string {
	return fmt.Sprintf("DateTimeTimeZone(DateTime: %s, TimeZone: %s)", s.DateTime.Format(time.RFC3339), *s.TimeZone)
}

// Equal compares the TimeSlot to the other value and returns true
// if the values are equal
func (s DateTimeTimeZone) Equal(other DateTimeTimeZone) bool {
	return s.DateTime.Equal(*other.DateTime) && s.TimeZone == other.TimeZone
}

func (s DateTimeTimeZone) Now() *DateTimeTimeZone {
	utc := time.Now().UTC()
	utcString := "(UTC) Coordinated Universal Time"
	return &DateTimeTimeZone{
		DateTime: &utc,
		TimeZone: &utcString,
	}
}

func (s DateTimeTimeZone) NowAdd(duration time.Duration) *DateTimeTimeZone {
	now := s.Now()
	nowAdd := now.DateTime.Add(duration)
	now.DateTime = &nowAdd
	return now
}

func (s *DateTimeTimeZone) MarshalJSON() ([]byte, error) {

	timezone, _ := globalSupportedTimeZones.GetTimeZoneByDisplayName(*s.TimeZone)

	stringVal := fmt.Sprintf("{ \"dateTime\": \"%v\", \"timeZone\": \"%v\" }",
		s.DateTime.Format("2006-01-02T15:04:05"), timezone)
	return []byte(stringVal), nil
}

// UnmarshalJSON implements the json unmarshal to be used by the json-library
func (s *DateTimeTimeZone) UnmarshalJSON(data []byte) error {
	tmp := struct {
		DateTime  *string  `json:"dateTime"`
		TimeZone  *string   `json:"timeZone"`
	}{}

	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}

	s.DateTime, err = parseTimeAndLocation(*tmp.DateTime, *tmp.TimeZone)
	if err != nil {
		return fmt.Errorf("cannot parse timestamp : %v", err)
	}
	s.TimeZone = tmp.TimeZone

	return nil
}

func parseTimeAndLocation(timeToParse, locationToParse string) (*time.Time, error) {
	parsedTime, err := time.Parse("2006-01-02T15:04:05.999999999", timeToParse)
	if err != nil {
		return nil, err
	}
	parsedTimeZone, err := time.LoadLocation(locationToParse)
	if err != nil {
		return nil, err
	}
	tz := parsedTime.In(parsedTimeZone)
	return &tz, nil
}

