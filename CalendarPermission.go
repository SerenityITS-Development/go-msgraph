package msgraph

import (
	"encoding/json"
	"fmt"
)

type CalendarPermission struct {
	ID                  	string
	IsRemovable	        	bool
	IsInsideOrganization	bool
	Role	        		string
	AllowedRoles	        []string
	EmailAppliedTo			EmailAddress

	calendar *Calendar
	graphClient *GraphClient // the graphClient that created this instance
}

func (cP CalendarPermission) String() string {
	return fmt.Sprintf("Calendar(ID: \"%v\", IsRemovable: \"%v\", IsInsideOrganization: \"%v\", Role: \"%v\","+
		"AllowedRoles: \"%v\", EmailAppliedTo: \"%v\")", cP.ID, cP.IsRemovable, cP.IsInsideOrganization, cP.Role,
		cP.AllowedRoles, cP.EmailAppliedTo)
}

func (cP *CalendarPermission) Delete(opts ...DeleteQueryOption) error {
	if cP.graphClient == nil {
		return ErrNotGraphClientSourced
	}

	resource := fmt.Sprintf("/users/%v/calendars/%v/calendarPermissions/%v", cP.calendar.Owner.Address,
		cP.calendar.ID, cP.ID)

	// TODO: check return body, maybe there is some potential success or error message hidden in it?
	// TODO: delete any child calendars
	err := cP.graphClient.makeDELETEAPICall(resource, compileDeleteQueryOptions(opts), nil)
	return err
}


// UnmarshalJSON implements the json unmarshal to be used by the json-library
func (cP *CalendarPermission) UnmarshalJSON(data []byte) error {
	tmp := struct {
		ID                  	string  		`json:"id"`
		IsRemovable	        	bool  			`json:"isRemovable"`
		IsInsideOrganization	bool  			`json:"isInsideOrganization"`
		Role	        		string  		`json:"role"`
		AllowedRoles	        []string		`json:"allowedRoles"`
		EmailAppliedTo			EmailAddress
	}{}
	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}

	cP.ID = tmp.ID
	cP.IsRemovable = tmp.IsRemovable
	cP.IsInsideOrganization = tmp.IsInsideOrganization
	cP.Role = tmp.Role
	cP.AllowedRoles = tmp.AllowedRoles
	cP.EmailAppliedTo = tmp.EmailAppliedTo

	return nil
}