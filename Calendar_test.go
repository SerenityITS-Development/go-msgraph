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

	permission, err := calendar.ShareReadWith(EmailAddress{Address: "taimana@outlook.com"}, false, false, "read")
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

	duration, err := time.ParseDuration("1h")
	if err != nil {
		log.Fatalf("failed to parse duration: %v", err)
	}

	attendees := Attendees{ {EmailAddress: EmailAddress{ Address: "taimana@outlook.com" }, Type: AttendeeRequired},
		{EmailAddress: EmailAddress{ Address: "taimana.nospam@outlook.com" }, Type: AttendeeOptional} }
	eventPost := CalendarEvent{
		Subject:               "Test",
		StartTime:             DateTimeTimeZone{}.Now(),
		EndTime:               DateTimeTimeZone{}.NowAdd(duration),
		Attendees:             &attendees,
		AllowNewTimeProposals: false,
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
		log.Fatalf("failed to create calendar: %v", err)
	}

	log.Println(time.Now())
	newEvent, err := calendar.CreateEvent(eventPost)
	if err != nil {
		log.Fatalf("failed to create event: %v", err)
	}
	log.Println(time.Now())

	err = newEvent.Delete()
	if err != nil {
		log.Fatalf("failed to delete event: %v", err)
	}
	log.Println(time.Now())

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