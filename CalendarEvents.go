package msgraph

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"
)

// CalendarEvents represents multiple events of a Calendar. The amount of entries is determined by the timespan that is used to load the Calendar
type CalendarEvents []CalendarEvent

func (c CalendarEvents) setGraphClient(gC *GraphClient) CalendarEvents {
	for i := range c {
		c[i] = c[i].setGraphClient(gC)
	}
	return c
}

func (c CalendarEvents) String() string {
	var events = make([]string, len(c))
	for i, calendarEvent := range c {
		events[i] = calendarEvent.String()
	}
	return fmt.Sprintf("CalendarEvents(%v)", strings.Join(events, ", "))
}

// PrettySimpleString returns all Calendar Events in a readable format, mostly used for logging purposes
func (c CalendarEvents) PrettySimpleString() string {
	var events = make([]string, len(c))
	for i, calendarEvent := range c {
		events[i] = calendarEvent.PrettySimpleString()
	}
	return fmt.Sprintf("CalendarEvents(%v)", strings.Join(events, ", "))
}

// SortByStartDateTime sorts the array in this CalendarEvents instance
func (c CalendarEvents) SortByStartDateTime() {
	sort.Slice(c, func(i, j int) bool {
		before := c[j].StartTime.DateTime
		return c[i].StartTime.DateTime.Before(*before) })
}

// GetCalendarEventsAtCertainTime returns a subset of CalendarEvents that either start or end
// at the givenTime or whose StartTime is before and EndTime is After the givenTime
func (c CalendarEvents) GetCalendarEventsAtCertainTime(givenTime time.Time) CalendarEvents {
	var events []CalendarEvent
	for _, event := range c {
		if event.StartTime.DateTime.Equal(givenTime) || event.EndTime.DateTime.Equal(givenTime) || (event.StartTime.DateTime.Before(givenTime) && event.EndTime.DateTime.After(givenTime)) {
			events = append(events, event)
		}
	}
	return events
}

func (c CalendarEvents) FindEventByTransactionId(value string) (*CalendarEvent, error) {
	for _, event := range c {
		if *event.TransactionID == value {
			return &event, nil
		}
	}
	return nil, ErrFindCalendarEvent
}

// Equal returns true if the two CalendarEvent[] are equal. The order of the events doesn't matter
func (c CalendarEvents) Equal(others CalendarEvents) bool {
Outer:
	for _, event := range c {
		for _, toCompare := range others {
			if event.Equal(toCompare) {
				continue Outer
			}
		}
		return false
	}
	return len(c) == len(others) // if we reach this, all CAlendarEvents in c have been found in others
}

// UnmarshalJSON implements the json unmarshal to be used by the json-library. The only
// purpose of this overwrite is to immediately sort the []CalendarEvent by StartDateTime
func (c *CalendarEvents) UnmarshalJSON(data []byte) error {
	tmp := struct {
		CalendarEvents []CalendarEvent `json:"value"`
	}{}

	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return fmt.Errorf("cannot UnmarshalJSON: %v | Data: %v", err, string(data))
	}

	*c = tmp.CalendarEvents // re-assign the

	c.SortByStartDateTime()
	return nil
}
