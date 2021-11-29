package msgraph

import (
	"encoding/json"
	"errors"
	"fmt"
)

type Phone	struct {
	Number		string  `json:"number,omitempty"`
	Type		PhoneType `json:"phoneType,omitempty"`
}

type PhoneType string

const (
	Home PhoneType 	= "home"
	Business		= "business"
	Mobile			= "mobile"
	Other			= "other"
	Assistant		= "assistant"
	HomeFax			= "homeFax"
	BusinessFax		= "businessFax"
	OtherFax		= "otherFax"
	Pager			= "pager"
	Radio			= "radio"
)



func (s Phone) String() string {
	return fmt.Sprintf("Phone(Number: %s, Type: %s)",
		s.Number, s.Type)
}

// Equal compares the OnlineMeetingInfo to the other value and returns true
// if the values are equal
func (s Phone) Equal(other Phone) bool {
	return s.Number == other.Number && s.Type == other.Type
}

// UnmarshalJSON implements the json unmarshal to be used by the json-library
func (s *Phone) UnmarshalJSON(data []byte) error {
	tmp := struct {
		Number		string    `json:"number"`
		Type	PhoneType `json:"phoneType"`
	}{}

	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}

	s.Number = tmp.Number
	s.Type = tmp.Type

	return nil
}

func (in PhoneType) IsValid() error {
	switch in {
	case Home, Business, Mobile, Other, Assistant, HomeFax, BusinessFax, OtherFax, Pager, Radio:
		return nil
	}
	return errors.New("invalid phone type")
}

func (in *PhoneType) UnmarshalJSON(data []byte) error {
	var s string
	//goland:noinspection GoUnhandledErrorResult
	json.Unmarshal(data, &s)
	out := PhoneType(s)
	switch out {
	case Home, Business, Mobile, Other, Assistant, HomeFax, BusinessFax, OtherFax, Pager, Radio:
		*in = out
		return nil
	}
	return errors.New("invalid phone type")
}
