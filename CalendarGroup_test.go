package msgraph

import (
	"log"
	"testing"
)

func TestCalendarGroup(t *testing.T) {
	t.Run("Create, List, GetByName, Delete CalendarGroup", func(t *testing.T) {
		client, err := NewGraphClient(
			msGraphTenantID,
			msGraphApplicationID,
			msGraphClientSecret)
		if err != nil {
			log.Fatalf("failed to create graph client: %v", err)
		}

		users, err := client.ListUsers()
		if err != nil {
			log.Fatalf("failed to list users: %v", err)
		}

		user, err := users.GetUserByMail(msGraphExistingUserPrincipalInGroup)
		if err != nil {
			log.Fatalf("failed to find user %s", msGraphExistingUserPrincipalInGroup)
		}

		calendarGroupName := "Test CG"
		calendarGroup, err := user.CreateCalendarGroup(calendarGroupName)
		if err != nil {
			log.Fatalf("failed to create calendar group: %v", err)
		}

		calendarGroups, err := user.ListCalendarGroups()
		if err != nil {
			log.Fatalf("failed to list calendar groups: %v", err)
		}

		calendarGroup, err = calendarGroups.GetByName(calendarGroupName)
		if err != nil {
			log.Fatalf("failed to get calendar group by name: %v", err)
		}

		err = calendarGroup.Delete()
		if err != nil {
			log.Fatalf("failed to delete calendar group: %v", err)
		}
	})

	t.Run("Get Calendar from My Calendars Group", func(t *testing.T) {
		client, err := NewGraphClient(
			msGraphTenantID,
			msGraphApplicationID,
			msGraphClientSecret)
		if err != nil {
			log.Fatalf("failed to create graph client: %v", err)
		}

		users, err := client.ListUsers()
		if err != nil {
			log.Fatalf("failed to list users: %v", err)
		}

		user, err := users.GetUserByMail(msGraphExistingUserPrincipalInGroup)
		if err != nil {
			log.Fatalf("failed to find user %s", msGraphExistingUserPrincipalInGroup)
		}

		calendarGroups, err := user.ListCalendarGroups()
		if err != nil {
			log.Fatalf("failed to list calendar groups: %v", err)
		}

		calendarGroupName := "My Calendars"
		calendarGroup, err := calendarGroups.GetByName(calendarGroupName)
		if err != nil {
			log.Fatalf("failed to get calendar group: %v", err)
		}

		calendarName := "Test Calendar"
		calendar, err := calendarGroup.CreateCalendar(calendarName)
		if err != nil {
			log.Fatalf("failed to get create calendar: %v", err)
		}

		err = calendar.Delete()
		if err != nil {
			log.Fatalf("failed to get delete calendar: %v", err)
		}
	})
}