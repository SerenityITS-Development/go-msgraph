package msgraph
//@file: goland:noinspection SpellCheckingInspection

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// CalendarEvent represents a single event within a calendar
type CalendarEvent struct {
	ID                    string  `json:"-"`
	Subject               string   `json:"subject,omitempty"`
	StartTime             DateTimeTimeZone   `json:"start,omitempty"`
	EndTime               DateTimeTimeZone    `json:"end,omitempty"` // endtime of the event, correct timezone is set
	CreatedDateTime       *time.Time    `json:"-"` // Creation time of the CalendarEvent, has the correct timezone set from OriginalStartTimeZone (json)
	LastModifiedDateTime  *time.Time       `json:"-"` // Last modified time of the CalendarEvent, has the correct timezone set from OriginalEndTimeZone (json)
	OriginalStartTimeZone *time.Location   `json:"-"` // The original start-timezone, is already integrated in the calendartimes. Caution: is UTC on full day events
	OriginalEndTimeZone   *time.Location   `json:"-"` // The original end-timezone, is already integrated in the calendartimes. Caution: is UTC on full day events
	ICalUID               string `json:"iCalUId,omitempty"`
	Importance            string    `json:"importance,omitempty"`
	Sensitivity           Sensitivity  `json:"sensitivity,omitempty"`
	IsAllDay              bool  `json:"isAllDay,omitempty"` // true = full day event, otherwise false
	IsCancelled           bool    `json:"isCancelled,omitempty"` // calendar event has been cancelled but is still in the calendar
	IsOrganizer           bool    `json:"isOrganizer,omitempty"` // true if the calendar owner is the organizer
	SeriesMasterID        string  `json:"seriesMasterId,omitempty"` // the ID of the master-entry of this series-event if any
	Type                  string `json:"type,omitempty"`
	ResponseStatus        *ResponseStatus `json:"responseStatus,omitempty"` // how the calendar-owner responded to the event (normally "organizer" because support-calendar is the host)

	Attendees      *Attendees  `json:"attendees,omitempty"` // represents all attendees to this CalendarEvent
	Organizer      *struct {
		EmailAddress EmailAddress `json:"emailAddress,omitempty"`
	} `json:"organizer,omitempty"`

	graphClient *GraphClient

	AllowNewTimeProposals 	bool `json:"allowNewTimeProposals,omitempty"`
	BodyPreview				string `json:"bodyPreview,omitempty"`
	Body					*struct {
		ContentType			string `json:"contentType,omitempty"`
		Content				string 	`json:"content,omitempty"`
	} `json:"body,omitempty"`
	Location				*CalendarLocation `json:"location,omitempty"`
	Locations				*[]CalendarLocation `json:"locations,omitempty"`
	HideAttendees			bool `json:"hideAttendees,omitempty"`
	CancelledOccurrences    []string `json:"cancelledOccurrences,omitempty"`
	Categories			    []string `json:"categories,omitempty"`
	ChangeKey				string `json:"-"`
	ExceptionOccurrences    []string `json:"-"`
	HasAttachments			bool `json:"hasAttachments,omitempty"`
	IsDraft					bool `json:"isDraft,omitempty"`
	IsOnlineMeeting			bool `json:"isOnlineMeeting,omitempty"`
	IsReminderOn			bool `json:"isReminderOn,omitempty"`
	OccurrenceID			string 	`json:"occurrenceId,omitempty"`
	OnlineMeetingInfo		*OnlineMeetingInfo `json:"onlineMeetingInfo,omitempty"`
	OnlineMeetingProvider	*OnlineMeetingProvider `json:"onlineMeetingProvider,omitempty"`
	OnlineMeetingURL		string `json:"onlineMeetingUrl,omitempty"`
	ReminderMinutesBeforeStart	int32  `json:"reminderMinutesBeforeStart,omitempty"`
	ResponseRequested		bool `json:"responseRequested,omitempty"`
	ShowAs					CalendarEventShowAs  `json:"showAs,omitempty"`
	TransactionID			string `json:"TransactionId,omitempty"`
	UUID					string `json:"uuid,omitempty"`
	WebLink					string `json:"webLink,omitempty"`
	Recurrence				*PatternedRecurrence `json:"recurrence,omitempty"`
}

// GetFirstAttendee returns the first Attendee that is not the organizer of the event from the Attendees array.
// If none is found then an Attendee with the Name of "None" will be returned.
func (c CalendarEvent) GetFirstAttendee() Attendee {
	for _, attendee := range *c.Attendees {
		if !attendee.EmailAddress.Equal(c.Organizer.EmailAddress) {
			return attendee
		}
	}

	return Attendee{ EmailAddress: EmailAddress{ Name: "None", Address: "None"}}
}

func (c CalendarEvent) String() string {
	return fmt.Sprintf("CalendarEvent(ID: \"%v\", CreatedDateTime: \"%v\", LastModifiedDateTime: \"%v\", "+
		"ICalUId: \"%v\", Subject: \"%v\", "+
		"Importance: \"%v\", Sensitivity: \"%v\", IsAllDay: \"%v\", IsCancelled: \"%v\", "+
		"IsOrganizer: \"%v\", SeriesMasterId: \"%v\", ShowAs: \"%v\", Type: \"%v\", ResponseStatus: \"%v\", "+
		"Attendees: \"%v\", Organizer: \"%v\", Start: \"%v\", End: \"%v\")", c.ID, c.CreatedDateTime, c.LastModifiedDateTime,
		c.ICalUID, c.Subject, c.Importance,
		c.Sensitivity, c.IsAllDay, c.IsCancelled, c.IsOrganizer, c.SeriesMasterID, c.ShowAs,
		c.Type, c.ResponseStatus, c.Attendees, c.Organizer.EmailAddress.Name+" "+c.Organizer.EmailAddress.Address,
		c.StartTime, c.EndTime)
}

// PrettySimpleString returns all Calendar Events in a readable format, mostly used for logging purposes
func (c CalendarEvent) PrettySimpleString() string {
	return fmt.Sprintf("{ %v (%v) [%v - %v] }", c.Subject, c.GetFirstAttendee().EmailAddress.Name, c.StartTime, c.EndTime)
}

func (c CalendarEvent) Delete(opts ...DeleteQueryOption) error {
	if c.graphClient == nil {
		return ErrNotGraphClientSourced
	}

	resource := fmt.Sprintf("/users/%v/events/%v", c.Organizer.EmailAddress.Address, c.ID)

	// TODO: check return body, maybe there is some potential success or error message hidden in it?
	err := c.graphClient.makeDELETEAPICall(resource, compileDeleteQueryOptions(opts), nil)
	return err
}

// Equal returns wether the CalendarEvent is identical to the given CalendarEvent
func (c CalendarEvent) Equal(other CalendarEvent) bool {
	return c.ID == other.ID && c.CreatedDateTime.Equal(*other.CreatedDateTime) && c.LastModifiedDateTime.Equal(*other.LastModifiedDateTime) &&
		c.ICalUID == other.ICalUID && c.Subject == other.Subject && c.Importance == other.Importance && c.Sensitivity == other.Sensitivity &&
		c.IsAllDay == other.IsAllDay && c.IsCancelled == other.IsCancelled && c.IsOrganizer == other.IsOrganizer &&
		c.SeriesMasterID == other.SeriesMasterID && c.ShowAs == other.ShowAs && c.Type == other.Type && c.ResponseStatus.Equal(*other.ResponseStatus) &&
		c.StartTime.Equal(other.StartTime) && c.EndTime.Equal(other.EndTime) &&
		c.Attendees.Equal(*other.Attendees) && c.Organizer.EmailAddress.Equal(other.Organizer.EmailAddress)
}

// UnmarshalJSON implements the json unmarshal to be used by the json-library
func (c *CalendarEvent) UnmarshalJSON(data []byte) error {
	tmp := struct {
		ID                    string            `json:"id"`
		CreatedDateTime       string            `json:"createdDateTime"`
		LastModifiedDateTime  string            `json:"lastModifiedDateTime"`
		OriginalStartTimeZone string            `json:"originalStartTimeZone"`
		OriginalEndTimeZone   string            `json:"originalEndTimeZone"`
		ICalUID               string            `json:"iCalUId"`
		Subject               string            `json:"subject"`
		Importance            string            `json:"importance"`
		Sensitivity           Sensitivity            `json:"sensitivity"`
		IsAllDay              bool              `json:"isAllDay"`
		IsCancelled           bool              `json:"isCancelled"`
		IsOrganizer           bool              `json:"isOrganizer"`
		SeriesMasterID        string            `json:"seriesMasterId"`
		ShowAs                CalendarEventShowAs            `json:"showAs"`
		Type                  string            `json:"type"`
		ResponseStatus        ResponseStatus    `json:"responseStatus"`
		//Start                 map[string]string `json:"start"`
		//End                   map[string]string `json:"end"`
		Start                 DateTimeTimeZone `json:"start"`
		End                   DateTimeTimeZone `json:"end"`
		Attendees             Attendees         `json:"attendees"`
		Organizer             struct {
			EmailAddress EmailAddress `json:"emailAddress"`
		} `json:"organizer"`

		AllowNewTimeProposals bool				`json:"allowNewTimeProposals"`
		HideAttendees			bool			`json:"hideAttendees"`
		BodyPreview				string			`json:"bodyPreview"`
		Body					struct {
			ContentType			string			`json:"contentType"`
			Content				string			`json:"content"`
		}			`json:"body"`
		CancelledOccurrences    []string		`json:"cancelledOccurrences"`
		Categories			    []string		`json:"categories"`
		ChangeKey				string			`json:"changeKey"`
		ExceptionOccurrences    []string		`json:"exceptionOccurrences"`
		HasAttachments			bool			`json:"hasAttachments"`
		IsDraft					bool			`json:"isDraft"`
		IsOnlineMeeting			bool			`json:"isOnlineMeeting"`
		IsReminderOn			bool			`json:"isReminderOn"`
		OccurrenceID			string			`json:"occurrenceId"`
		OnlineMeetingInfo		OnlineMeetingInfo	`json:"onlineMeetingInfo"`
		OnlineMeetingProvider	OnlineMeetingProvider `json:"onlineMeetingProvider"`
		OnlineMeetingURL			string			`json:"onlineMeetingUrl"`
		ReminderMinutesBeforeStart	int32  `json:"reminderMinutesBeforeStart"`
		ResponseRequested		bool  `json:"responseRequested"`
		TransactionID			string  `json:"TransactionId"`
		UUID					string  `json:"uuid"`
		WebLink					string  `json:"webLink"`
		Recurrence				PatternedRecurrence `json:"recurrence"`
		Location				CalendarLocation `json:"location"`
		Locations				[]CalendarLocation `json:"locations"`
	}{}


	var err error
	// unmarshal to tmp-struct, return if error
	if err = json.Unmarshal(data, &tmp); err != nil {
		return fmt.Errorf("error on json.Unmarshal: %v | Data: %v", err, string(data))
	}

	c.ID = tmp.ID

	created, err := time.Parse(time.RFC3339Nano, tmp.CreatedDateTime)
	if err != nil {
		return fmt.Errorf("cannot time.Parse with RFC3339Nano createdDateTime %v: %v", tmp.CreatedDateTime, err)
	}
	modified, err := time.Parse(time.RFC3339Nano, tmp.LastModifiedDateTime)
	if err != nil {
		return fmt.Errorf("cannot time.Parse with RFC3339Nano lastModifiedDateTime %v: %v", tmp.LastModifiedDateTime, err)
	}

	c.CreatedDateTime = &created
	c.LastModifiedDateTime = &modified

	c.OriginalStartTimeZone, err = mapTimeZoneStrings(tmp.OriginalStartTimeZone)
	if err != nil {
		return fmt.Errorf("cannot time.LoadLocation originalStartTimeZone %v: %v", tmp.OriginalStartTimeZone, err)
	}
	c.OriginalEndTimeZone, err = mapTimeZoneStrings(tmp.OriginalEndTimeZone)
	if err != nil {
		return fmt.Errorf("cannot time.LoadLocation originalEndTimeZone %v: %v", tmp.OriginalEndTimeZone, err)
	}
	c.ICalUID = tmp.ICalUID
	c.Subject = tmp.Subject
	c.Importance = tmp.Importance
	c.Sensitivity = tmp.Sensitivity
	c.IsAllDay = tmp.IsAllDay
	c.IsCancelled = tmp.IsCancelled
	c.IsOrganizer = tmp.IsOrganizer
	c.SeriesMasterID = tmp.SeriesMasterID
	c.ShowAs = tmp.ShowAs
	c.Type = tmp.Type
	c.ResponseStatus = &tmp.ResponseStatus
	if tmp.Attendees != nil {
		c.Attendees = &tmp.Attendees
	}


	c.Organizer = (*struct {
		EmailAddress EmailAddress `json:"emailAddress,omitempty"`
	})(&tmp.Organizer)
	c.AllowNewTimeProposals = tmp.AllowNewTimeProposals
	c.BodyPreview = tmp.BodyPreview
	c.HideAttendees = tmp.HideAttendees
	c.Body = (*struct {
		ContentType string `json:"contentType,omitempty"`
		Content     string `json:"content,omitempty"`
	})(&tmp.Body)
	/*
	c.StartTime, err = time.Parse(time.RFC3339Nano, tmp.Start)
	if err != nil {
		return fmt.Errorf("cannot parse timestamp with RFC3339Nano: %v", err)
	}
	c.EndTime, err = time.Parse(time.RFC3339Nano, tmp.End)
	if err != nil {
		return fmt.Errorf("cannot parse timestamp with RFC3339Nano: %v", err)
	}

	 */
	c.CancelledOccurrences = tmp.CancelledOccurrences
	c.Categories = tmp.Categories
	c.ChangeKey = tmp.ChangeKey
	c.ExceptionOccurrences = tmp.ExceptionOccurrences
	c.HasAttachments = tmp.HasAttachments
	c.IsDraft = tmp.IsDraft
	c.IsOnlineMeeting = tmp.IsOnlineMeeting
	c.IsReminderOn = tmp.IsReminderOn
	c.OccurrenceID = tmp.OccurrenceID
	c.OnlineMeetingInfo = &tmp.OnlineMeetingInfo
	c.OnlineMeetingProvider = &tmp.OnlineMeetingProvider
	c.OnlineMeetingURL = tmp.OnlineMeetingURL
	c.ReminderMinutesBeforeStart = tmp.ReminderMinutesBeforeStart
	c.ResponseRequested = tmp.ResponseRequested
	c.ShowAs = tmp.ShowAs
	c.TransactionID = tmp.TransactionID
	c.UUID = tmp.UUID
	c.WebLink = tmp.WebLink
	c.Recurrence = &tmp.Recurrence
	c.Locations = &tmp.Locations
	c.Location = &tmp.Location
	c.StartTime = tmp.Start
	c.EndTime = tmp.End


	/*
	// Parse event start & endtime with timezone
	c.StartTime, err = parseTimeAndLocation(tmp.Start["dateTime"], tmp.Start["timeZone"]) // the timeZone is normally ALWAYS UTC, microsoft converts time date & time to that
	if err != nil {
		return fmt.Errorf("cannot parse start-dateTime %v AND timeZone %v: %v", tmp.Start["dateTime"], tmp.Start["timeZone"], err)
	}
	c.EndTime, err = parseTimeAndLocation(tmp.End["dateTime"], tmp.End["timeZone"]) // the timeZone is normally ALWAYS UTC, microsoft converts time date & time to that
	if err != nil {
		return fmt.Errorf("cannot parse end-dateTime %v AND timeZone %v: %v", tmp.End["dateTime"], tmp.End["timeZone"], err)
	}

	 */

	// Hint: OriginalStartTimeZone & end are UTC (set by microsoft) if it is a full-day event, this will be handled in the next section
	c.StartTime.DateTime = c.StartTime.DateTime.In(c.OriginalStartTimeZone) // move the StartTime to the orignal start-timezone
	c.StartTime.TimeZone = c.OriginalStartTimeZone.String()
	c.EndTime.DateTime = c.EndTime.DateTime.In(c.OriginalEndTimeZone)       // move the EndTime to the orignal end-timezone
	c.EndTime.TimeZone = c.OriginalEndTimeZone.String()

	// Now check if it's a full-day event, if yes, the event is UTC anyway. We need it to be accurate for the program to work
	// hence we set it to time.Local. It can later be manipulated by the program to a different timezone but the times also have
	// to be recalculated. E.g. we set it to UTC+2 hence it will start at 02:00 and end at 02:00, not 00:00 -> manually set to 00:00
	if c.IsAllDay && FullDayEventTimeZone != time.UTC {
		// set to local location
		c.StartTime.DateTime = c.StartTime.DateTime.In(FullDayEventTimeZone)
		c.StartTime.TimeZone = FullDayEventTimeZone.String()
		c.EndTime.DateTime = c.EndTime.DateTime.In(FullDayEventTimeZone)
		c.EndTime.TimeZone = FullDayEventTimeZone.String()
		// get offset in seconds
		_, startOffSet := c.StartTime.DateTime.Zone()
		_, endOffSet := c.EndTime.DateTime.Zone()
		// decrease time to 00:00 again
		c.StartTime.DateTime = c.StartTime.DateTime.Add(-1 * time.Second * time.Duration(startOffSet))
		c.EndTime.DateTime = c.EndTime.DateTime.Add(-1 * time.Second * time.Duration(endOffSet))
	}

	return nil
}

/*
// parseTimeAndLocation is just a helper method to shorten the code in the Unmarshal json
func parseTimeAndLocation(timeToParse, locationToParse string) (time.Time, error) {
	parsedTime, err := time.Parse("2006-01-02T15:04:05.999999999", timeToParse)
	if err != nil {
		return time.Time{}, err
	}
	parsedTimeZone, err := time.LoadLocation(locationToParse)
	if err != nil {
		return time.Time{}, err
	}
	return parsedTime.In(parsedTimeZone), nil
}

 */


// mapTimeZoneStrings maps various Timezones used by Microsoft to go-understandable timezones or returns the source-zone if no mapping is found
func mapTimeZoneStrings(timeZone string) (*time.Location, error) {
	if timeZone == "tzone://Microsoft/Custom" {
		return FullDayEventTimeZone, nil
	}
	tz, err := globalSupportedTimeZones.GetTimeZoneByAlias(timeZone)
	if err == nil {
		return tz, nil
	}
	location, err := time.LoadLocation(timeZone)
	if err == nil {
		return location, nil
	}
	return nil, errors.New("oh-oh")
}
