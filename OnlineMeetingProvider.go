package msgraph

import (
	"encoding/json"
	"errors"
)

type OnlineMeetingProvider string

const (
	OnlineProviderUnknown OnlineMeetingProvider 	= "unknown"
	TeamsBusiness		= "teamsForBusiness"
	SkypeBusiness			= "skypeForBusiness"
	Skype			= "skypeForConsumer"
)


func (in OnlineMeetingProvider) IsValid() error {
	switch in {
	case OnlineProviderUnknown, TeamsBusiness, SkypeBusiness, Skype:
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
	case OnlineProviderUnknown, TeamsBusiness, SkypeBusiness, Skype:
		*in = out
		return nil
	}
	return errors.New("invalid provider type")
}
