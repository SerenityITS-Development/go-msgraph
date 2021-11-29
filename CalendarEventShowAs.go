package msgraph

import (
	"encoding/json"
	"errors"
)

type CalendarEventShowAs string

//goland:noinspection SpellCheckingInspection
const (
	Free CalendarEventShowAs 	= "free"
	Tentative		= "tentative"
	Busy			= "busy"
	OOTO			= "oof"
	WorkingElsewhere = "workingElsewhere"
	ShowCalendarAsUnknown			= "unknown"
)


func (in CalendarEventShowAs) IsValid() error {
	switch in {
	case ShowCalendarAsUnknown, Free, Tentative, Busy, OOTO, WorkingElsewhere:
		return nil
	}
	return errors.New("invalid show as type")
}

func (in *CalendarEventShowAs) UnmarshalJSON(data []byte) error {
	var s string
	//goland:noinspection GoUnhandledErrorResult
	json.Unmarshal(data, &s)
	out := CalendarEventShowAs(s)
	switch out {
	case ShowCalendarAsUnknown, Free, Tentative, Busy, OOTO, WorkingElsewhere:
		*in = out
		return nil
	}
	return errors.New("invalid show as type")
}
