package msgraph

import (
	"encoding/json"
	"fmt"
)

// CalendarLocation represents a location for CalendarEvent
type CalendarLocation struct {
	DisplayName				string `json:"displayName,omitempty"`
	LocationType			string `json:"locationType,omitempty"`
	LocationURI				string `json:"locationUri,omitempty"`
	LocationEmailAddress	string `json:"LocationEmailAddress,omitempty"`
	Address					*Address `json:"address,omitempty"`
	Coordinates 			*OutlookCoordinates `json:"coordinates,omitempty"`
}

func (s CalendarLocation) String() string {
	return fmt.Sprintf("CalendarLocation(DisplayName: %s, LocationType: %s, LocationURI: %s, LocationEmailAddress: %s, Address: %s, Coordinates: %s)",
		s.DisplayName, s.LocationType, s.LocationURI, s.LocationEmailAddress, s.Address.String(), s.Coordinates.String())
}

// Equal compares the CalendarLocation to the other value and returns true
// if the values are equal
func (s CalendarLocation) Equal(other CalendarLocation) bool {
	return s.DisplayName == other.DisplayName && s.LocationURI == other.LocationURI && s.LocationType == other.LocationType &&
		s.LocationEmailAddress == other.LocationEmailAddress && s.Address.Equal(*other.Address) && s.Coordinates.Equal(*other.Coordinates)
}

// UnmarshalJSON implements the json unmarshal to be used by the json-library
func (s *CalendarLocation) UnmarshalJSON(data []byte) error {
	tmp := struct {
		DisplayName				string `json:"displayName,omitempty"`
		LocationType			string `json:"locationType,omitempty"`
		LocationURI				string `json:"locationUri,omitempty"`
		LocationEmailAddress	string `json:"LocationEmailAddress,omitempty"`
		Address					*Address `json:"address,omitempty"`
		Coordinates 			*OutlookCoordinates `json:"coordinates,omitempty"`
	}{}

	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}

	s.DisplayName = tmp.DisplayName
	s.LocationType = tmp.LocationType
	s.LocationURI = tmp.LocationURI
	s.LocationEmailAddress = tmp.LocationEmailAddress
	s.Address = tmp.Address
	s.Coordinates = tmp.Coordinates

	return nil
}
