package msgraph

import (
	"encoding/json"
	"errors"
)

type Sensitivity string

//goland:noinspection SpellCheckingInspection
const (
	SensitivityNormal   Sensitivity = "normal"
	SensitivityPersonal = "personal"
	SensitivityPrivate      = "private"
	SensitivityConfidential = "confidential"
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
	OnlineProviderTeamsBusiness = "teamsForBusiness"
	OnlineProviderSkypeBusiness = "skypeForBusiness"
	OnlineProviderSkype         = "skypeForConsumer"
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
	ImportanceNormal		 = "normal"
	ImportanceHigh			 = "high"
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
	ContentTypeHtml		 = "html"
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