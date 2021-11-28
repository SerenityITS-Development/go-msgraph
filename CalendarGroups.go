package msgraph

import (
	"strings"
)

type CalendarGroups []CalendarGroup

func (c CalendarGroups) String() string {
	var calendarGroups = make([]string, len(c))
	for i, calendarGroup := range c {
		calendarGroups[i] = calendarGroup.String()
	}
	return "Calendars(" + strings.Join(calendarGroups, " | ") + ")"
}

// setGraphClient sets the graphClient instance in this instance and all child-instances (if any)
func (c CalendarGroups) setGraphClient(gC *GraphClient, u *User) CalendarGroups {
	for i := range c {
		c[i].setGraphAndUser(gC, u)
	}
	return c
}

// GetByName Gets the calendar by name
// Supports optional OData query parameters https://docs.microsoft.com/en-us/graph/query-parameters
func (c CalendarGroups) GetByName(name string) (CalendarGroup, error) {
	for _, calendarGroup := range c {
		if calendarGroup.Name == name {
			return calendarGroup, nil
		}
	}
	return CalendarGroup{}, ErrFindCalendarGroup
}