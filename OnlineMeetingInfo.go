package msgraph

import (
	"encoding/json"
	"fmt"
	"strings"
)

// OnlineMeetingInfo represents a proposed time for CalendarEvent
type OnlineMeetingInfo struct {
	ConferenceID	string `json:"conferenceId,omitempty"`
	JoinURL			string    `json:"joinUrl,omitempty"`
	QuickDial		string   `json:"quickDial,omitempty"`
	TollFreeNumbers	[]string  `json:"tollFreeNumbers,omitempty"`
	TollNumber		string `json:"tollNumber,omitempty"`
	Phones			*[]Phone `json:"phones,omitempty"`
}

func (s OnlineMeetingInfo) String() string {
	return fmt.Sprintf("OnlineMeetingInfo(ConferenceID: %s, JoinURL: %s, QuickDial: %s, TollNumber: %s, " +
		" TollFreeNumbers: %s)",
		s.ConferenceID, s.JoinURL, s.QuickDial, s.TollNumber, strings.Join(s.TollFreeNumbers, "|"))
}

// Equal compares the OnlineMeetingInfo to the other value and returns true
// if the values are equal
func (s OnlineMeetingInfo) Equal(other OnlineMeetingInfo) bool {
	return s.ConferenceID == other.ConferenceID
}

// UnmarshalJSON implements the json unmarshal to be used by the json-library
func (s *OnlineMeetingInfo) UnmarshalJSON(data []byte) error {
	tmp := struct {

		ConferenceID	string    `json:"conferenceId"`
		JoinURL			string    `json:"joinUrl"`
		QuickDial		string    `json:"quickDial"`
		TollFreeNumbers	[]string  `json:"tollFreeNumbers"`
		TollNumber		string    `json:"tollNumber"`
		Phones			[]Phone	  `json:"phones"`
	}{}

	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}

	s.ConferenceID = tmp.ConferenceID
	s.JoinURL = tmp.JoinURL
	s.QuickDial = tmp.QuickDial
	s.TollNumber = tmp.TollNumber
	s.TollFreeNumbers = tmp.TollFreeNumbers
	s.Phones = &tmp.Phones

	return nil
}
