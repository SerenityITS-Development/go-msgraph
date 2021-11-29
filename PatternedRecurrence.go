package msgraph

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type DayOfWeek string

const (
	Sunday DayOfWeek = "sunday"
	Monday			= "monday"
	Tuesday			= "tuesday"
	Wednesday		= "wednesday"
	Thursday		= "thursday"
	Friday			= "friday"
	Saturday		= "saturday"
)

func (in DayOfWeek) IsValid() error {
	switch in {
	case Sunday, Monday, Tuesday, Wednesday, Thursday, Friday, Saturday:
		return nil
	}
	return errors.New("invalid day of week type")
}

func (in *DayOfWeek) UnmarshalJSON(data []byte) error {
	var s string
	//goland:noinspection GoUnhandledErrorResult
	json.Unmarshal(data, &s)
	out := DayOfWeek(s)
	switch out {
	case Sunday, Monday, Tuesday, Wednesday, Thursday, Friday, Saturday:
		*in = out
		return nil
	}
	return errors.New("invalid day of week type")
}

func DayOfWeekArrayEquals(left, right []DayOfWeek) bool {
	if len(left) != len(right) {
		return false
	}
	for i, v := range left {
		if v != right[i] {
			return false
		}
	}
	return true
}

type WeekIndex string

const (
	FirstWeek WeekIndex = "first"
	SecondWeek			= "second"
	ThirdWeek			 = "third"
	ForthWeek			= "forth"
	LastWeek			= "last"
)


func (in WeekIndex) IsValid() error {
	switch in {
	case FirstWeek, SecondWeek, ThirdWeek, ForthWeek, LastWeek:
		return nil
	}
	return errors.New("invalid week index type")
}

func (in *WeekIndex) UnmarshalJSON(data []byte) error {
	var s string
	//goland:noinspection GoUnhandledErrorResult
	json.Unmarshal(data, &s)
	out := WeekIndex(s)
	switch out {
	case FirstWeek, SecondWeek, ThirdWeek, ForthWeek, LastWeek:
		*in = out
		return nil
	}
	return errors.New("invalid week index type")
}

type RecurrencePatternType string

const (
	Daily RecurrencePatternType = "daily"
	Weekly						= "weekly"
	AbsoluteMonthly				= "absoluteMonthly"
	RelativeMonthly				= "relativeMonthly"
	AbsoluteYearly				= "absoluteYearly"
	RelativeYearly				= "relativeYearly"
)

func (in RecurrencePatternType) IsValid() error {
	switch in {
	case Daily, Weekly, AbsoluteMonthly, RelativeMonthly, AbsoluteYearly, RelativeYearly:
		return nil
	}
	return errors.New("invalid recurrence pattern type")
}

func (in *RecurrencePatternType) UnmarshalJSON(data []byte) error {
	var s string
	//goland:noinspection GoUnhandledErrorResult
	json.Unmarshal(data, &s)
	out := RecurrencePatternType(s)
	switch out {
	case Daily, Weekly, AbsoluteMonthly, RelativeMonthly, AbsoluteYearly, RelativeYearly:
		*in = out
		return nil
	}
	return errors.New("invalid recurrence pattern type")
}

type RecurrenceRangeType string

const (
	RecurEndDate	RecurrenceRangeType = "endDate"
	RecurNoEnd							= "noEnd"
	RecurNumbered						= "numbered"
)

func (in RecurrenceRangeType) IsValid() error {
	switch in {
	case RecurEndDate, RecurNoEnd, RecurNumbered:
		return nil
	}
	return errors.New("invalid recur type")
}

func (in *RecurrenceRangeType) UnmarshalJSON(data []byte) error {
	var s string
	//goland:noinspection GoUnhandledErrorResult
	json.Unmarshal(data, &s)
	out := RecurrenceRangeType(s)
	switch out {
	case RecurEndDate, RecurNoEnd, RecurNumbered:
		*in = out
		return nil
	}
	return errors.New("invalid recur type")
}

type RecurrenceRange struct {
	EndDate			time.Time
	NumberOfOccurrences	int32
	RecurrenceTimeZone string
	StartDate		time.Time
	Type			RecurrenceRangeType
}

func (s RecurrenceRange) String() string {
	return fmt.Sprintf("RecurrenceRange(StartDate: %s, EndDate: %s, NumberOfOccurrences: %v, RecurrenceTimeZone: %s, Type: %s)",
		s.StartDate.Format(time.RFC3339Nano), s.EndDate.Format(time.RFC3339Nano), s.NumberOfOccurrences,
		s.RecurrenceTimeZone, s.Type)
}

func (s RecurrenceRange) Equal(other RecurrenceRange) bool {
	return s.StartDate.Equal(other.StartDate) && s.EndDate.Equal(other.EndDate) &&
		s.NumberOfOccurrences == other.NumberOfOccurrences && s.RecurrenceTimeZone == other.RecurrenceTimeZone &&
		s.Type == other.Type
}

func (s *RecurrenceRange) UnmarshalJSON(data []byte) error {
	tmp := struct {
		EndDate			string   `json:"endDate"`
		NumberOfOccurrences	int32   `json:"numberOfOccurrences"`
		RecurrenceTimeZone string   `json:"recurrenceTimeZone"`
		StartDate		string   `json:"startDate"`
		Type			RecurrenceRangeType   `json:"type"`
	}{}

	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}

	s.StartDate, err = time.Parse(time.RFC3339Nano, tmp.StartDate)
	if err != nil {
		return fmt.Errorf("cannot parse timestamp with RFC3339Nano: %v", err)
	}
	s.EndDate, err = time.Parse(time.RFC3339Nano, tmp.EndDate) // the timeZone is normally ALWAYS UTC, microsoft converts time date & time to that, but it does not matter here
	if err != nil {
		return fmt.Errorf("cannot parse timestamp with RFC3339Nano: %v", err)
	}
	s.NumberOfOccurrences = tmp.NumberOfOccurrences
	s.RecurrenceTimeZone = tmp.RecurrenceTimeZone
	s.Type = tmp.Type

	return nil
}

type RecurrencePattern struct {
	DayOfMonth		int32
	DaysOfWeek		[]DayOfWeek
	FirstDayOfWeek	DayOfWeek
	Index			WeekIndex
	Interval		int32
	Month			int32
	Type			RecurrencePatternType
}

func (s RecurrencePattern) String() string {

	return fmt.Sprintf("RecurrencePattern(DayOfMonth: %v, FirstDayOfWeek: %v, Index: %v, Interval: %v, Month: %v, Type: %v, DaysOfWeeks: %s)",
		s.DayOfMonth, s.FirstDayOfWeek, s.Index, s.Interval, s.Month, s.Type, "TODO")
}

func (s RecurrencePattern) Equal(other RecurrencePattern) bool {
	return s.DayOfMonth == other.DayOfMonth && DayOfWeekArrayEquals(s.DaysOfWeek, other.DaysOfWeek) &&
		s.FirstDayOfWeek == other.FirstDayOfWeek && s.Index == other.Index && s.Interval == other.Interval &&
		s.Month == other.Month && s.Type == other.Type
}

func (s *RecurrencePattern) UnmarshalJSON(data []byte) error {
	tmp := struct {
		DayOfMonth		int32   `json:"endDate"`
		DaysOfWeek		[]DayOfWeek   `json:"daysOfWeek"`
		FirstDayOfWeek	DayOfWeek   `json:"firstDayOfWeek"`
		Index			WeekIndex   `json:"index"`
		Interval		int32   `json:"interval"`
		Month			int32   `json:"month"`
		Type			RecurrencePatternType   `json:"type"`
	}{}

	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}

	s.DayOfMonth = tmp.DayOfMonth
	s.DaysOfWeek = tmp.DaysOfWeek
	s.FirstDayOfWeek = tmp.FirstDayOfWeek
	s.Index = tmp.Index
	s.Interval = tmp.Interval
	s.Month = tmp.Month
	s.Type = tmp.Type

	return nil
}

type PatternedRecurrence struct {
	Pattern    	*RecurrencePattern `json:"pattern,omitempty"`
	Range 		*RecurrenceRange `json:"range,omitempty"`
}


func (s PatternedRecurrence) String() string {

	return fmt.Sprintf("RecurrencePattern(Pattern: %v, RecurrenceRange: %v)",
		s.Pattern.String(), s.Range.String())
}

func (s PatternedRecurrence) Equal(other PatternedRecurrence) bool {
	return s.Pattern.Equal(*other.Pattern) && s.Range.Equal(*other.Range)
}

func (s *PatternedRecurrence) UnmarshalJSON(data []byte) error {
	tmp := struct {
		Pattern    	RecurrencePattern `json:"pattern"`
		Range 		RecurrenceRange `json:"range"`
	}{}

	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}

	s.Pattern = &tmp.Pattern
	s.Range = &tmp.Range

	return nil
}
