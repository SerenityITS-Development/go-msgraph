package msgraph

import (
	"encoding/json"
	"errors"
)

type Sensitivity string

//goland:noinspection SpellCheckingInspection
const (
	Normal Sensitivity 	= "normal"
	Personal		= "personal"
	Private			= "private"
	Confidential			= "confidential"
)


func (in Sensitivity) IsValid() error {
	switch in {
	case Normal, Personal, Private, Confidential:
		return nil
	}
	return errors.New("invalid sensitivity type")
}

func (in *Sensitivity) UnmarshalJSON(data []byte) error {
	var s string
	//goland:noinspection GoUnhandledErrorResult
	json.Unmarshal(data, &s)
	out := Sensitivity(s)
	switch out {
	case Normal, Personal, Private, Confidential:
		*in = out
		return nil
	}
	return errors.New("invalid sensitivity type")
}

