package msgraph

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

// Calendar represents a single calendar of a user
//
// See https://developer.microsoft.com/en-us/graph/docs/api-reference/v1.0/resources/calendar
type Calendar struct {
	ID                  string // The group's unique identifier. Read-only.
	Name                string // The calendar name.
	CanEdit             bool   // True if the user can write to the calendar, false otherwise. This property is true for the user who created the calendar. This property is also true for a user who has been shared a calendar and granted write access.
	CanShare            bool   // True if the user has the permission to share the calendar, false otherwise. Only the user who created the calendar can share it.
	CanViewPrivateItems bool   // True if the user can read calendar items that have been marked private, false otherwise.
	ChangeKey           string // Identifies the version of the calendar object. Every time the calendar is changed, changeKey changes as well. This allows Exchange to apply changes to the correct version of the object. Read-only.

	Owner EmailAddress // If set, this represents the user who created or added the calendar. For a calendar that the user created or added, the owner property is set to the user. For a calendar shared with the user, the owner property is set to the person who shared that calendar with the user.

	graphClient *GraphClient // the graphClient that created this instance
}

func (c Calendar) String() string {
	return fmt.Sprintf("Calendar(ID: \"%v\", Name: \"%v\", canEdit: \"%v\", canShare: \"%v\", canViewPrivateItems: \"%v\", ChangeKey: \"%v\", "+
		"Owner: \"%v\")", c.ID, c.Name, c.CanEdit, c.CanShare, c.CanViewPrivateItems, c.ChangeKey, c.Owner)
}

// setGraphClient sets the graphClient instance in this instance and all child-instances (if any)
func (c *Calendar) setGraphClient(graphClient *GraphClient) {
	c.graphClient = graphClient
	c.Owner.setGraphClient(graphClient)
}

func (c *Calendar) ShareReadWith(email EmailAddress, isInsideOrganization bool,
		isRemovable bool, role string, opts ...CreateQueryOption) (CalendarPermission, error) {

	if c.graphClient == nil {
		return CalendarPermission{}, ErrNotGraphClientSourced
	}

	resource := fmt.Sprintf("/users/%v/calendars/%v/calendarPermissions", c.Owner.Address, c.ID)

	calendarPermission := CalendarPermission{graphClient: c.graphClient, calendar: c}
	bodyBytes, err := json.Marshal(struct {
		IsInsideOrganization bool `json:"isInsideOrganization"`
		IsRemovable bool `json:"isRemovable"`
		Role string `json:"role"`
		EmailAppliedTo EmailAddress `json:"emailAddress"`
	}{
		IsRemovable: isRemovable,
		IsInsideOrganization: isInsideOrganization,
		Role: role,
		EmailAppliedTo: email,
	})
	if err != nil {
		return calendarPermission, err
	}

	reader := bytes.NewReader(bodyBytes)
	err = c.graphClient.makePOSTAPICall(resource, compileCreateQueryOptions(opts), reader, &calendarPermission)
	return calendarPermission, err
}

func (c Calendar) CreateEvent(event CalendarEvent, opts ...CreateQueryOption) (*CalendarEvent, error) {
	if c.graphClient == nil {
		return nil, ErrNotGraphClientSourced
	}

	user, err := c.graphClient.GetUser(c.Owner.Address)
	if err != nil {
		return nil, err
	}

	if len(globalSupportedTimeZones.Value) == 0 {
		var err error
		// TODO: this is a dirty fix, because opts could contain other things than a context, e.g. select
		// parameters. This could produce unexpected outputs and therefore break the globalSupportedTimeZones variable.
		globalSupportedTimeZones, err = user.getTimeZoneChoices(compileCreateQueryOptions(opts))
		if err != nil {
			return nil, err
		}
	}

	resource := fmt.Sprintf("/users/%v/calendars/%v/events", c.Owner.Address, c.ID)
	newEvent := CalendarEvent{graphClient: c.graphClient}
	bodyBytes, err := json.Marshal(event)
	if err != nil {
		return &newEvent, err
	}


	reader := bytes.NewReader(bodyBytes)
	err = c.graphClient.makePOSTAPICall(resource, compileCreateQueryOptions(opts), reader, &newEvent)

	return &newEvent, err
}

// Delete deletes this calendar instance for this user. Use with caution.
//
// Reference: https://docs.microsoft.com/en-us/graph/api/user-delete
func (c Calendar) Delete(opts ...DeleteQueryOption) error {
	if c.graphClient == nil {
		return ErrNotGraphClientSourced
	}

	resource := fmt.Sprintf("/users/%v/calendars/%v", c.Owner.Address, c.ID)

	// TODO: check return body, maybe there is some potential success or error message hidden in it?
	err := c.graphClient.makeDELETEAPICall(resource, compileDeleteQueryOptions(opts), nil)
	return err
}

func (c Calendar) ListEvents(startDateTime, endDateTime time.Time, opts ...ListQueryOption) (CalendarEvents, error) {
	if c.graphClient == nil {
		return CalendarEvents{}, ErrNotGraphClientSourced
	}

	user, err := c.graphClient.GetUser(c.Owner.Address)
	if err != nil {
		return nil, err
	}

	if len(globalSupportedTimeZones.Value) == 0 {
		var err error
		// TODO: this is a dirty fix, because opts could contain other things than a context, e.g. select
		// parameters. This could produce unexpected outputs and therefore break the globalSupportedTimeZones variable.
		globalSupportedTimeZones, err = user.getTimeZoneChoices(compileListQueryOptions(opts))
		if err != nil {
			return CalendarEvents{}, err
		}
	}

	resource := fmt.Sprintf("/users/%v/calendars/%v/events", c.Owner.Address, c.ID)

	// set GET-Params for start and end time
	var reqOpt = compileListQueryOptions(opts)
	reqOpt.queryValues.Add("startdatetime", startDateTime.Format("2006-01-02T00:00:00"))
	reqOpt.queryValues.Add("enddatetime", endDateTime.Format("2006-01-02T00:00:00"))

	var calendarEvents CalendarEvents
	err = c.graphClient.makeGETAPICall(resource, reqOpt, &calendarEvents)
	if err != nil {
		return nil, err
	}
	calendarEvents.setGraphClient(c.graphClient)
	return calendarEvents, nil
}

// UnmarshalJSON implements the json unmarshal to be used by the json-library
func (c *Calendar) UnmarshalJSON(data []byte) error {
	tmp := struct {
		ID                  string `json:"id"`                  // the calendars ID
		Name                string `json:"name"`                // the name of the calendar
		CanShare            bool   `json:"canShare"`            // true if the current account can shares this calendar
		CanViewPrivateItems bool   `json:"canViewPrivateItems"` // true if the current account can view private entries
		CanEdit             bool   `json:"canEdit"`             // true if the current account can edit the calendar
		ChangeKey           string `json:"changeKey"`
		Owner               EmailAddress
	}{}
	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}

	c.ID = tmp.ID
	c.Name = tmp.Name
	c.CanEdit = tmp.CanEdit
	c.CanShare = tmp.CanShare
	c.CanViewPrivateItems = tmp.CanViewPrivateItems
	c.ChangeKey = tmp.ChangeKey

	c.Owner = tmp.Owner

	return nil
}
