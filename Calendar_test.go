package msgraph

import (
	"fmt"
	"log"
	"testing"
	"time"
)

func TestCalendar_String(t *testing.T) {
	testCalendars := GetTestListCalendars(t)
	if skipCalendarTests {
		t.Skip("Skipping due to missing 'MSGraphExistingCalendarsOfUser' value")
	}

	for _, testCalendar := range testCalendars {
		tt := struct {
			name string
			c    *Calendar
			want string
		}{
			name: "Test first calendar",
			c:    &testCalendar,
			want: fmt.Sprintf("Calendar(ID: \"%v\", Name: \"%v\", canEdit: \"%v\", canShare: \"%v\", canViewPrivateItems: \"%v\", ChangeKey: \"%v\", "+
				"Owner: \"%v\")", testCalendar.ID, testCalendar.Name, testCalendar.CanEdit, testCalendar.CanShare, testCalendar.CanViewPrivateItems, testCalendar.ChangeKey, testCalendar.Owner),
		}
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.String(); got != tt.want {
				t.Errorf("Calendar.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalendar_ShareReadWith(t *testing.T) {
	client, err := NewGraphClient(
		msGraphTenantID,
		msGraphApplicationID,
		msGraphClientSecret)
	if err != nil {
		log.Fatalf("failed to create graph client: %v", err)
	}

	user, err := client.GetUser(msGraphExistingUserPrincipalInGroup)
	if err != nil {
		log.Fatalf("failed to get user: %v", err)
	}

	group, err := user.CreateCalendarGroup("Shared")
	if err != nil {
		log.Fatalf("failed to create calendar group: %v", err)
	}

	calendar, err := group.CreateCalendar("Shareable")
	if err != nil {
		log.Fatalf("failed to create calendar: %v", err)
	}

	permission, err := calendar.ShareReadWith(EmailAddress{Address: "doesnoexst_email@outlook.com"}, false, false, "read")
	if err != nil {
		log.Fatalf("failed to share calendar: %v", err)
	}

	err = permission.Delete()
	if err != nil {
		log.Fatalf("failed to delete permission: %v", err)
	}

	err = calendar.Delete()
	if err != nil {
		log.Fatalf("failed to delete calendar: %v", err)
	}

	err = group.Delete()
	if err != nil {
		log.Fatalf("failed to delete calendar group: %v", err)
	}
}



func TestCalendar_CreateEvent(t *testing.T) {

	duration := 1 * time.Hour
	testString := "Test"
	nowTime := DateTimeTimeZone{}.Now()
	endTime := DateTimeTimeZone{}.NowAdd(duration)
	falseValue := false
	transactionId := "1"
	required := AttendeeRequired
	optional := AttendeeOptional

	attendees := Attendees{ {EmailAddress: &EmailAddress{ Address: "doesnotexist_email1@outlook.com" }, Type: &required},
		{EmailAddress: &EmailAddress{ Address: "doesnotexist_email2@outlook.com" }, Type: &optional} }
	eventPost := CalendarEvent{
		Subject:               &testString,
		StartTime:             nowTime,
		EndTime:               endTime,
		Attendees:             &attendees,
		AllowNewTimeProposals: &falseValue,
		TransactionID:         &transactionId,
	}

	client, err := NewGraphClient(
		msGraphTenantID,
		msGraphApplicationID,
		msGraphClientSecret)
	if err != nil {
		log.Fatalf("failed to create graph client: %v", err)
	}

	user, err := client.GetUser(msGraphExistingUserPrincipalInGroup)
	if err != nil {
		log.Fatalf("failed to get user: %v", err)
	}

	group, err := user.CreateCalendarGroup("Events")
	if err != nil {
		log.Fatalf("failed to create calendar group: %v", err)
	}

	calendar, err := group.CreateCalendar("Event Calendar")
	if err != nil {
		err2 := err
		if group.ID != "" {
			err = group.Delete()
			if err != nil {
				log.Fatalf("failed to delete calendar group: %v", err)
			}
		}
		log.Fatalf("failed to create calendar: %v", err2)
	}

	newEvent, err := calendar.CreateEvent(eventPost)
	if err != nil {
		err2 := err
		if calendar.ID != "" {
			err = calendar.Delete()
			if err != nil {
				log.Fatalf("failed to delete calendar: %v", err)
			}
		}

		if group.ID != "" {
			err = group.Delete()
			if err != nil {
				log.Fatalf("failed to delete calendar group: %v", err)
			}
		}
		log.Fatalf("failed to create event: %v", err2)
	}

	endDuration := 10 * time.Hour
	startDuration := -10 * time.Hour
	eventsList, err := calendar.ListEvents(time.Now().Add(startDuration), time.Now().Add(endDuration))
	if err != nil {
		err2 := err
		if calendar.ID != "" {
			err = calendar.Delete()
			if err != nil {
				log.Fatalf("failed to delete calendar: %v", err)
			}
		}

		if group.ID != "" {
			err = group.Delete()
			if err != nil {
				log.Fatalf("failed to delete calendar group: %v", err)
			}
		}
		log.Fatalf("failed to list events: %v", err2)
	}

	_, err = eventsList.FindEventByTransactionId("1")
	if err != nil {
		err2 := err
		if calendar.ID != "" {
			err = calendar.Delete()
			if err != nil {
				log.Fatalf("failed to delete calendar: %v", err)
			}
		}

		if group.ID != "" {
			err = group.Delete()
			if err != nil {
				log.Fatalf("failed to delete calendar group: %v", err)
			}
		}
		log.Fatalf("failed to find event: %v", err2)
	}

	err = newEvent.Update()
	if err != nil {
		err2 := err
		if calendar.ID != "" {
			err = calendar.Delete()
			if err != nil {
				log.Fatalf("failed to delete calendar: %v", err)
			}
		}

		if group.ID != "" {
			err = group.Delete()
			if err != nil {
				log.Fatalf("failed to delete calendar group: %v", err)
			}
		}
		log.Fatalf("failed to update event: %v", err2)
	}

	err = newEvent.Delete()
	if err != nil {
		err2 := err
		if calendar.ID != "" {
			err = calendar.Delete()
			if err != nil {
				log.Fatalf("failed to delete calendar: %v", err)
			}
		}

		if group.ID != "" {
			err = group.Delete()
			if err != nil {
				log.Fatalf("failed to delete calendar group: %v", err)
			}
		}
		log.Fatalf("failed to delete event: %v", err2)
	}

	if calendar.ID != "" {
		err = calendar.Delete()
		if err != nil {
			log.Fatalf("failed to delete calendar: %v", err)
		}
	}

	if group.ID != "" {
		err = group.Delete()
		if err != nil {
			log.Fatalf("failed to delete calendar group: %v", err)
		}
	}
}