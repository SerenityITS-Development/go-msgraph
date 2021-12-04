package msgraph

import (
	"encoding/json"
	"errors"
)

type CalendarEventShowAs string

//goland:noinspection SpellCheckingInspection
const (
	ShowAsFree      CalendarEventShowAs = "free"
	ShowAsTentative = "tentative"
	ShowAsBusy       = "busy"
	ShowAsOOTO              = "oof"
	ShowAssWorkingElsewhere = "workingElsewhere"
	ShowAsUnknown           = "unknown"
)


func (in CalendarEventShowAs) IsValid() error {
	switch in {
	case ShowAsUnknown, ShowAsFree, ShowAsTentative, ShowAsBusy, ShowAsOOTO, ShowAssWorkingElsewhere:
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
	case ShowAsUnknown, ShowAsFree, ShowAsTentative, ShowAsBusy, ShowAsOOTO, ShowAssWorkingElsewhere:
		*in = out
		return nil
	}
	return errors.New("invalid show as type")
}
