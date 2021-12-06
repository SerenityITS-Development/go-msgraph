package msgraph

import (
	"encoding/json"
	"fmt"
)

// Attendee struct represents an attendee for a CalendarEvent
type Attendee struct {
	Type           *AttendeeType   `json:"type,omitempty"`
	ResponseStatus *ResponseStatus `json:"status,omitempty"`
	ProposedNewTime *TimeSlot `json:"proposedNewTime,omitempty"`
	EmailAddress    *EmailAddress `json:"emailAddress,omitempty"`
}

func (a Attendee) String() string {
	return fmt.Sprintf("Type: %v, E-mail: %s, ResponseStatus: %v, TimeSlot: %s",
		a.Type, a.EmailAddress.String(), a.ResponseStatus, a.ProposedNewTime.String())
}

// Equal compares the Attendee to the other Attendee and returns true
// if the two given Attendees are equal, Otherwise returns false
func (a Attendee) Equal(other Attendee) bool {
	return a.Type == other.Type && a.EmailAddress.Address == other.EmailAddress.Address && a.ResponseStatus.Equal(*other.ResponseStatus)
}

// UnmarshalJSON implements the json unmarshal to be used by the json-library
func (a *Attendee) UnmarshalJSON(data []byte) error {
	tmp := struct {
		Type         *AttendeeType         `json:"type,omitempty"`
		Status       *ResponseStatus `json:"status,omitempty"`
		EmailAddress    *EmailAddress `json:"emailAddress,omitempty"`
		ProposedNewTime *TimeSlot     `json:"proposedNewTime,omitempty"`
	}{}

	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return fmt.Errorf("attendee: %v", err.Error())
	}

	a.Type = tmp.Type
	a.EmailAddress = tmp.EmailAddress
	a.ResponseStatus = tmp.Status
	a.ProposedNewTime = tmp.ProposedNewTime

	return nil
}
