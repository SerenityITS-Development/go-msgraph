package msgraph
//@file: goland:noinspection SpellCheckingInspection

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// CalendarEvent represents a single event within a calendar
type CalendarEvent struct {
	ID                    *string  `json:"id,omitempty"`
	Subject               *string   `json:"subject,omitempty"`
	StartTime             *DateTimeTimeZone   `json:"start,omitempty"`
	EndTime               *DateTimeTimeZone    `json:"end,omitempty"` // endtime of the event, correct timezone is set
	CreatedDateTime       *time.Time    `json:"-"` // Creation time of the CalendarEvent, has the correct timezone set from OriginalStartTimeZone (json)
	LastModifiedDateTime  *time.Time       `json:"-"` // Last modified time of the CalendarEvent, has the correct timezone set from OriginalEndTimeZone (json)
	OriginalStartTimeZone *time.Location   `json:"-"` // The original start-timezone, is already integrated in the calendartimes. Caution: is UTC on full day events
	OriginalEndTimeZone   *time.Location   `json:"-"` // The original end-timezone, is already integrated in the calendartimes. Caution: is UTC on full day events
	ICalUID               *string `json:"iCalUId,omitempty"`
	Importance            *Importance    `json:"importance,omitempty"`
	Sensitivity           *Sensitivity  `json:"sensitivity,omitempty"`
	IsAllDay              *bool  `json:"isAllDay,omitempty"` // true = full day event, otherwise false
	IsCancelled           *bool    `json:"isCancelled,omitempty"` // calendar event has been cancelled but is still in the calendar
	IsOrganizer           *bool    `json:"isOrganizer,omitempty"` // true if the calendar owner is the organizer
	SeriesMasterID        *string  `json:"seriesMasterId,omitempty"` // the ID of the master-entry of this series-event if any
	Type                  *string `json:"type,omitempty"`
	ResponseStatus        *ResponseStatus `json:"responseStatus,omitempty"` // how the calendar-owner responded to the event (normally "organizer" because support-calendar is the host)

	Attendees      *Attendees  `json:"attendees,omitempty"` // represents all attendees to this CalendarEvent
	Organizer      *Organizer `json:"organizer,omitempty"`

	graphClient *GraphClient

	AllowNewTimeProposals 	*bool `json:"allowNewTimeProposals,omitempty"`
	BodyPreview				*string `json:"bodyPreview,omitempty"`
	Body					*ContentBody `json:"body,omitempty"`
	Location				*CalendarLocation `json:"location,omitempty"`
	Locations				*[]CalendarLocation `json:"locations,omitempty"`
	HideAttendees			*bool `json:"hideAttendees,omitempty"`
	CancelledOccurrences    *[]string `json:"cancelledOccurrences,omitempty"`
	Categories			    *[]string `json:"categories,omitempty"`
	ChangeKey				*string `json:"-"`
	ExceptionOccurrences    *[]string `json:"-"`
	HasAttachments			*bool `json:"hasAttachments,omitempty"`
	IsDraft					*bool `json:"isDraft,omitempty"`
	IsOnlineMeeting			*bool `json:"isOnlineMeeting,omitempty"`
	IsReminderOn			*bool `json:"isReminderOn,omitempty"`
	OccurrenceID			*string 	`json:"occurrenceId,omitempty"`
	OnlineMeetingInfo		*OnlineMeetingInfo `json:"onlineMeetingInfo,omitempty"`
	OnlineMeetingProvider	*OnlineMeetingProvider `json:"onlineMeetingProvider,omitempty"`
	OnlineMeetingURL		*string `json:"onlineMeetingUrl,omitempty"`
	ReminderMinutesBeforeStart	*int32  `json:"reminderMinutesBeforeStart,omitempty"`
	ResponseRequested		*bool `json:"responseRequested,omitempty"`
	ShowAs					*CalendarEventShowAs  `json:"showAs,omitempty"`
	TransactionID			*string `json:"TransactionId,omitempty"`
	UUID					*string `json:"uuid,omitempty"`
	WebLink					*string `json:"webLink,omitempty"`
	Recurrence				*PatternedRecurrence `json:"recurrence,omitempty"`
}

func (c CalendarEvent) setGraphClient(gC *GraphClient) CalendarEvent {
	c.graphClient = gC
	return c
}

// GetFirstAttendee returns the first Attendee that is not the organizer of the event from the Attendees array.
// If none is found then an Attendee with the Name of "None" will be returned.
func (c CalendarEvent) GetFirstAttendee() Attendee {
	for _, attendee := range *c.Attendees {
		if !attendee.EmailAddress.Equal(*c.Organizer.EmailAddress) {
			return attendee
		}
	}

	return Attendee{ EmailAddress: &EmailAddress{ Name: "None", Address: "None"}}
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

func (c *CalendarEvent) Update(opts ... UpdateQueryOption) error {
	if c.graphClient == nil {
		return ErrNotGraphClientSourced
	}

	resource := fmt.Sprintf("/users/%v/events/%v", c.Organizer.EmailAddress.Address, *c.ID)

	bodyBytes, err := json.Marshal(c)
	if err != nil {
		return err
	}
	reader := bytes.NewReader(bodyBytes)

	// TODO: check return body, maybe there is some potential success or error message hidden in it?
	err = c.graphClient.makePATCHAPICall(resource, compileUpdateQueryOptions(opts), reader, nil)
	return err
}

func (c CalendarEvent) Delete(opts ...DeleteQueryOption) error {
	if c.graphClient == nil {
		return ErrNotGraphClientSourced
	}

	resource := fmt.Sprintf("/users/%v/events/%v", c.Organizer.EmailAddress.Address, *c.ID)

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
		c.StartTime.Equal(*other.StartTime) && c.EndTime.Equal(*other.EndTime) &&
		c.Attendees.Equal(*other.Attendees) && c.Organizer.EmailAddress.Equal(*other.Organizer.EmailAddress)
}

// UnmarshalJSON implements the json unmarshal to be used by the json-library
func (c *CalendarEvent) UnmarshalJSON(data []byte) error {
	tmp := struct {
		ID                    *string            `json:"id,omitempty"`
		CreatedDateTime       *string            `json:"createdDateTime,omitempty"`
		LastModifiedDateTime  *string            `json:"lastModifiedDateTime,omitempty"`
		OriginalStartTimeZone *string            `json:"originalStartTimeZone,omitempty"`
		OriginalEndTimeZone   *string            `json:"originalEndTimeZone,omitempty"`
		ICalUID               *string            `json:"iCalUId,omitempty"`
		Subject               *string            `json:"subject,omitempty"`
		Importance            *Importance            `json:"importance,omitempty"`
		Sensitivity           *Sensitivity            `json:"sensitivity,omitempty"`
		IsAllDay              *bool              `json:"isAllDay,omitempty"`
		IsCancelled           *bool              `json:"isCancelled,omitempty"`
		IsOrganizer           *bool              `json:"isOrganizer,omitempty"`
		SeriesMasterID        *string            `json:"seriesMasterId,omitempty"`
		ShowAs                *CalendarEventShowAs            `json:"showAs,omitempty"`
		Type                  *string            `json:"type,omitempty"`
		ResponseStatus        *ResponseStatus    `json:"responseStatus,omitempty"`
		Start                 *DateTimeTimeZone `json:"start,omitempty"`
		End                   *DateTimeTimeZone `json:"end,omitempty"`
		Attendees             *Attendees         `json:"attendees,omitempty"`
		Organizer             *Organizer `json:"organizer,omitempty"`

		AllowNewTimeProposals *bool				`json:"allowNewTimeProposals,omitempty"`
		HideAttendees			*bool			`json:"hideAttendees,omitempty"`
		BodyPreview				*string			`json:"bodyPreview,omitempty"`
		Body					*ContentBody			`json:"body,omitempty"`
		CancelledOccurrences    *[]string		`json:"cancelledOccurrences,omitempty"`
		Categories			    *[]string		`json:"categories,omitempty"`
		ChangeKey				*string			`json:"changeKey,omitempty"`
		ExceptionOccurrences    *[]string		`json:"exceptionOccurrences,omitempty"`
		HasAttachments			*bool			`json:"hasAttachments,omitempty"`
		IsDraft					*bool			`json:"isDraft,omitempty"`
		IsOnlineMeeting			*bool			`json:"isOnlineMeeting,omitempty"`
		IsReminderOn			*bool			`json:"isReminderOn,omitempty"`
		OccurrenceID			*string			`json:"occurrenceId,omitempty"`
		OnlineMeetingInfo		*OnlineMeetingInfo	`json:"onlineMeetingInfo,omitempty"`
		OnlineMeetingProvider	*OnlineMeetingProvider `json:"onlineMeetingProvider,omitempty"`
		OnlineMeetingURL			*string			`json:"onlineMeetingUrl,omitempty"`
		ReminderMinutesBeforeStart	*int32  `json:"reminderMinutesBeforeStart,omitempty"`
		ResponseRequested		*bool  `json:"responseRequested,omitempty"`
		TransactionID			*string  `json:"TransactionId,omitempty"`
		UUID					*string  `json:"uuid,omitempty"`
		WebLink					*string  `json:"webLink,omitempty"`
		Recurrence				*PatternedRecurrence `json:"recurrence,omitempty"`
		Location				*CalendarLocation `json:"location,omitempty"`
		Locations				*[]CalendarLocation `json:"locations,omitempty"`
	}{}


	var err error
	// unmarshal to tmp-struct, return if error
	if err = json.Unmarshal(data, &tmp); err != nil {
		return fmt.Errorf("error on json.Unmarshal: %v | Data: %v", err, string(data))
	}

	c.ID = tmp.ID

	created, err := time.Parse(time.RFC3339Nano, *tmp.CreatedDateTime)
	if err != nil {
		return fmt.Errorf("cannot time.Parse with RFC3339Nano createdDateTime %v: %v", tmp.CreatedDateTime, err)
	}
	modified, err := time.Parse(time.RFC3339Nano, *tmp.LastModifiedDateTime)
	if err != nil {
		return fmt.Errorf("cannot time.Parse with RFC3339Nano lastModifiedDateTime %v: %v", tmp.LastModifiedDateTime, err)
	}

	c.CreatedDateTime = &created
	c.LastModifiedDateTime = &modified

	c.OriginalStartTimeZone, err = mapTimeZoneStrings(*tmp.OriginalStartTimeZone)
	if err != nil {
		return fmt.Errorf("cannot time.LoadLocation originalStartTimeZone %v: %v", tmp.OriginalStartTimeZone, err)
	}
	c.OriginalEndTimeZone, err = mapTimeZoneStrings(*tmp.OriginalEndTimeZone)
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
	c.ResponseStatus = tmp.ResponseStatus
	if tmp.Attendees != nil {
		c.Attendees = tmp.Attendees
	}

	c.Organizer = tmp.Organizer
	c.AllowNewTimeProposals = tmp.AllowNewTimeProposals
	c.BodyPreview = tmp.BodyPreview
	c.HideAttendees = tmp.HideAttendees
	c.Body = tmp.Body
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
	c.OnlineMeetingInfo = tmp.OnlineMeetingInfo
	c.OnlineMeetingProvider = tmp.OnlineMeetingProvider
	c.OnlineMeetingURL = tmp.OnlineMeetingURL
	c.ReminderMinutesBeforeStart = tmp.ReminderMinutesBeforeStart
	c.ResponseRequested = tmp.ResponseRequested
	c.ShowAs = tmp.ShowAs
	c.TransactionID = tmp.TransactionID
	c.UUID = tmp.UUID
	c.WebLink = tmp.WebLink
	c.Recurrence = tmp.Recurrence
	c.Locations = tmp.Locations
	c.Location = tmp.Location
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

	originalStartTimeIn := c.StartTime.DateTime.In(c.OriginalStartTimeZone)
	originalEndTimeIn := c.EndTime.DateTime.In(c.OriginalEndTimeZone)
	originalStartTimeZoneString := c.OriginalStartTimeZone.String()
	originalEndTimeZoneString := c.OriginalEndTimeZone.String()
	// Hint: OriginalStartTimeZone & end are UTC (set by microsoft) if it is a full-day event, this will be handled in the next section
	c.StartTime.DateTime = &originalStartTimeIn // move the StartTime to the orignal start-timezone
	c.StartTime.TimeZone = &originalStartTimeZoneString
	c.EndTime.DateTime = &originalEndTimeIn      // move the EndTime to the orignal end-timezone
	c.EndTime.TimeZone = &originalEndTimeZoneString

	// Now check if it's a full-day event, if yes, the event is UTC anyway. We need it to be accurate for the program to work
	// hence we set it to time.Local. It can later be manipulated by the program to a different timezone but the times also have
	// to be recalculated. E.g. we set it to UTC+2 hence it will start at 02:00 and end at 02:00, not 00:00 -> manually set to 00:00
	if *c.IsAllDay && FullDayEventTimeZone != time.UTC {
		// set to local location
		fullStartTimeIn := c.StartTime.DateTime.In(FullDayEventTimeZone)
		fullEndTimeIn := c.EndTime.DateTime.In(FullDayEventTimeZone)
		fullStartTimeZoneString := FullDayEventTimeZone.String()
		fullEndTimeZoneString := FullDayEventTimeZone.String()
		c.StartTime.DateTime = &fullStartTimeIn
		c.StartTime.TimeZone = &fullStartTimeZoneString
		c.EndTime.DateTime = &fullEndTimeIn
		c.EndTime.TimeZone = &fullEndTimeZoneString
		// get offset in seconds
		_, startOffSet := c.StartTime.DateTime.Zone()
		_, endOffSet := c.EndTime.DateTime.Zone()
		// decrease time to 00:00 again
		fullStartTimeAdd := c.StartTime.DateTime.Add(-1 * time.Second * time.Duration(startOffSet))
		fullEndTimeAdd := c.EndTime.DateTime.Add(-1 * time.Second * time.Duration(endOffSet))
		c.StartTime.DateTime = &fullStartTimeAdd
		c.EndTime.DateTime = &fullEndTimeAdd
	}

	return nil
}

type ContentBody struct {
	ContentType			*ContentType			`json:"contentType,omitempty"`
	Content				*string			`json:"content,omitempty"`
}

type Organizer struct {
	EmailAddress *EmailAddress `json:"emailAddress,omitempty"`
}

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
