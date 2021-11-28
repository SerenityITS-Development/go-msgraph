package msgraph

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type CalendarGroup struct {
	Name      string
	ClassID   string
	ChangeKey string
	ID        string

	user *User

	graphClient *GraphClient
}

func (cG CalendarGroup) String() string {
	return fmt.Sprintf("CalendarGroup(ID: \"%v\", Name: \"%v\", classID: \"%v\", ChangeKey: \"%v\")",
		cG.ID, cG.Name, cG.ClassID, cG.ChangeKey)
}

func (cG *CalendarGroup) setGraphClient(graphClient *GraphClient) {
	cG.graphClient = graphClient
}
// CreateCalendar creates calendar inside the calendar group.
// Supports optional OData query parameters https://docs.microsoft.com/en-us/graph/query-parameters
//
// Reference: https://developer.microsoft.com/en-us/graph/docs/api-reference/v1.0/api/user-post-calendars
func (cG *CalendarGroup) CreateCalendar(name string, opts ...CreateQueryOption) (Calendar, error) {

	if cG.graphClient == nil {
		return Calendar{}, ErrNotGraphClientSourced
	}

	resource := fmt.Sprintf("/users/%v/calendarGroups/%v/calendars", cG.user.ID, cG.ID)
	calendar := Calendar{graphClient: cG.graphClient, calendarGroup: cG}
	bodyBytes, err := json.Marshal(struct {
		Name string `json:"name"`
	}{Name: name})
	if err != nil {
		return Calendar{}, err
	}
	if err != nil {
		return calendar, err
	}

	reader := bytes.NewReader(bodyBytes)
	err = cG.graphClient.makePOSTAPICall(resource, compileCreateQueryOptions(opts), reader, &calendar)

	calendar.setGraphClient(cG.graphClient)
	return calendar, err
}

// ListCalendars returns all calendars associated to that user and group.
// Supports optional OData query parameters https://docs.microsoft.com/en-us/graph/query-parameters
//
// Reference: https://developer.microsoft.com/en-us/graph/docs/api-reference/v1.0/api/user-list-calendars
func (cG CalendarGroup) ListCalendars(opts ...ListQueryOption) (Calendars, error) {
	if cG.graphClient == nil {
		return Calendars{}, ErrNotGraphClientSourced
	}
	resource := fmt.Sprintf("/users/%v/calendarGroups/%v/calendars", cG.user.ID, cG.ID)

	var marsh struct {
		Calendars Calendars `json:"value"`
	}
	err := cG.graphClient.makeGETAPICall(resource, compileListQueryOptions(opts), &marsh)
	marsh.Calendars.setGraphClient(cG.graphClient)
	return marsh.Calendars, err
}

// Delete deletes this user instance at the Microsoft Azure AD. Use with caution.
//
// Reference: https://docs.microsoft.com/en-us/graph/api/user-delete
func (cG CalendarGroup) Delete(opts ...DeleteQueryOption) error {
	if cG.graphClient == nil {
		return ErrNotGraphClientSourced
	}

	resource := fmt.Sprintf("/users/%v/calendarGroups/%v", cG.user.ID, cG.ID)

	// TODO: check return body, maybe there is some potential success or error message hidden in it?
	// TODO: delete any child calendars
	err := cG.graphClient.makeDELETEAPICall(resource, compileDeleteQueryOptions(opts), nil)
	return err
}


// UnmarshalJSON implements the json unmarshal to be used by the json-library
func (cG *CalendarGroup) UnmarshalJSON(data []byte) error {
	tmp := struct {
		Name      string `json:"name"`
		ClassID   string `json:"classId"`
		ChangeKey string `json:"changeKey"`
		ID        string `json:"id"`
	}{}
	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}

	cG.ID = tmp.ID
	cG.Name = tmp.Name
	cG.ClassID = tmp.ClassID
	cG.ChangeKey = tmp.ChangeKey

	return nil
}