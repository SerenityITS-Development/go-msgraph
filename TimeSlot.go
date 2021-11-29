package msgraph

import (
	"encoding/json"
	"fmt"
	"time"
)

// TimeSlot represents a proposed time for CalendarEvent
type TimeSlot struct {
	StartTime time.Time  `json:"start,omitempty"`   // status of the response, may be organizer, accepted, declined etc.
	EndTime   time.Time    `json:"end,omitempty"`// represents the time when the response was performed
}

func (s TimeSlot) String() string {
	return fmt.Sprintf("TimeSlot(StartTime: %s, EndTime: %s)", s.StartTime.Format(time.RFC3339Nano), s.EndTime.Format(time.RFC3339Nano))
}

// Equal compares the TimeSlot to the other value and returns true
// if the values are equal
func (s TimeSlot) Equal(other TimeSlot) bool {
	return s.StartTime.Equal(other.StartTime) && s.EndTime.Equal(other.EndTime)
}

// UnmarshalJSON implements the json unmarshal to be used by the json-library
func (s *TimeSlot) UnmarshalJSON(data []byte) error {
	tmp := struct {
		StartTime string   `json:"start"`
		EndTime   string   `json:"end"`
	}{}

	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}

	s.StartTime, err = time.Parse(time.RFC3339Nano, tmp.StartTime)
	if err != nil {
		return fmt.Errorf("cannot parse timestamp with RFC3339Nano: %v", err)
	}
	s.EndTime, err = time.Parse(time.RFC3339Nano, tmp.EndTime) // the timeZone is normally ALWAYS UTC, microsoft converts time date & time to that, but it does not matter here
	if err != nil {
		return fmt.Errorf("cannot parse timestamp with RFC3339Nano: %v", err)
	}

	return nil
}
