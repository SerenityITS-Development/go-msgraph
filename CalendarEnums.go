package msgraph

import (
	"encoding/json"
	"errors"
)


type EventResponseType string

//goland:noinspection SpellCheckingInspection
const (
	EventResponseNone   EventResponseType = "none"
	EventResponseOrganizer EventResponseType = "organizer"
	EventResponseTentativelyAccepted    EventResponseType  = "tentativelyAccepted"
	EventResponseAccepted EventResponseType = "accepted"
	EventResponseDeclined EventResponseType = "declined"
	EventResponseNotResponded EventResponseType = "notResponded"
)


func (in EventResponseType) IsValid() error {
	switch in {
	case EventResponseDeclined, EventResponseNone, EventResponseOrganizer,
		EventResponseAccepted, EventResponseTentativelyAccepted, EventResponseNotResponded:
		return nil
	}
	return errors.New("invalid EventResponseType type")
}

func (in *EventResponseType) UnmarshalJSON(data []byte) error {
	var s string
	//goland:noinspection GoUnhandledErrorResult
	json.Unmarshal(data, &s)
	out := EventResponseType(s)
	switch out {
	case EventResponseDeclined, EventResponseNone, EventResponseOrganizer,
		EventResponseAccepted, EventResponseTentativelyAccepted, EventResponseNotResponded:
		*in = out
		return nil
	}
	return errors.New("invalid EventResponseType type")
}

type Sensitivity string

//goland:noinspection SpellCheckingInspection
const (
	SensitivityNormal   Sensitivity = "normal"
	SensitivityPersonal Sensitivity = "personal"
	SensitivityPrivate    Sensitivity  = "private"
	SensitivityConfidential Sensitivity = "confidential"
)


func (in Sensitivity) IsValid() error {
	switch in {
	case SensitivityNormal, SensitivityPersonal, SensitivityPrivate, SensitivityConfidential:
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
	case SensitivityNormal, SensitivityPersonal, SensitivityPrivate, SensitivityConfidential:
		*in = out
		return nil
	}
	return errors.New("invalid sensitivity type")
}


type OnlineMeetingProvider string

const (
	OnlineProviderUnknown       OnlineMeetingProvider = "unknown"
	OnlineProviderTeamsBusiness OnlineMeetingProvider = "teamsForBusiness"
	OnlineProviderSkypeBusiness OnlineMeetingProvider = "skypeForBusiness"
	OnlineProviderSkype         OnlineMeetingProvider = "skypeForConsumer"
)


func (in OnlineMeetingProvider) IsValid() error {
	switch in {
	case OnlineProviderUnknown, OnlineProviderTeamsBusiness, OnlineProviderSkypeBusiness, OnlineProviderSkype:
		return nil
	}
	return errors.New("invalid provider type")
}

func (in *OnlineMeetingProvider) UnmarshalJSON(data []byte) error {
	var s string
	//goland:noinspection GoUnhandledErrorResult
	json.Unmarshal(data, &s)
	out := OnlineMeetingProvider(s)
	switch out {
	case OnlineProviderUnknown, OnlineProviderTeamsBusiness, OnlineProviderSkypeBusiness, OnlineProviderSkype:
		*in = out
		return nil
	}
	return errors.New("invalid provider type")
}

type Importance string

const (
	ImportanceLow Importance = "low"
	ImportanceNormal Importance 		 = "normal"
	ImportanceHigh	Importance 		 = "high"
)

func (in Importance) IsValid() error {
	switch in {
	case ImportanceLow, ImportanceNormal, ImportanceHigh:
		return nil
	}
	return errors.New("invalid importance type")
}

func (in *Importance) UnmarshalJSON(data []byte) error {
	var s string
	//goland:noinspection GoUnhandledErrorResult
	json.Unmarshal(data, &s)
	out := Importance(s)
	switch out {
	case ImportanceLow, ImportanceNormal, ImportanceHigh:
		*in = out
		return nil
	}
	return errors.New("invalid importance type")
}

type ContentType string

const (
	ContentTypeText ContentType = "text"
	ContentTypeHtml	ContentType	 = "html"
)

func (in ContentType) IsValid() error {
	switch in {
	case ContentTypeText, ContentTypeHtml:
		return nil
	}
	return errors.New("invalid ContentType type")
}

func (in *ContentType) UnmarshalJSON(data []byte) error {
	var s string
	//goland:noinspection GoUnhandledErrorResult
	json.Unmarshal(data, &s)
	out := ContentType(s)
	switch out {
	case ContentTypeText, ContentTypeHtml:
		*in = out
		return nil
	}
	return errors.New("invalid ContentType type")
}