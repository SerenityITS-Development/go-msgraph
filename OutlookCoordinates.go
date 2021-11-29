package msgraph

import (
	"encoding/json"
	"fmt"
)

// TimeSlot represents a proposed time for CalendarEvent
type OutlookCoordinates struct {
	Accuracy    float64 `json:"accuracy,omitempty"`
	Altitude    float64 `json:"altitude,omitempty"`
	AltitudeAccuracy float64 `json:"altitudeAccuracy,omitempty"`
	Latitude	float64`json:"latitude,omitempty"`
	Longitude	float64 `json:"longitude,omitempty"`
}

func (s OutlookCoordinates) String() string {
	return fmt.Sprintf("OutlookCoordinates(Accuracy: %v, Altitude: %v, AltitudeAccuracy: %v, Latitude: %v, Longitude: %v)",
		s.Accuracy, s.Altitude, s.AltitudeAccuracy, s.Latitude, s.Longitude)
}

var EPSILON float64 = 0.00000001
func floatEquals(a, b float64) bool {
	if (a - b) < EPSILON && (b - a) < EPSILON {
		return true
	}
	return false
}

// Equal compares the TimeSlot to the other value and returns true
// if the values are equal
func (s OutlookCoordinates) Equal(other OutlookCoordinates) bool {
	return floatEquals(s.Altitude, other.Altitude) && floatEquals(s.AltitudeAccuracy, other.AltitudeAccuracy) &&
		floatEquals(s.Accuracy, other.Accuracy) && floatEquals(s.Latitude, other.Latitude) &&
		floatEquals(s.Longitude, other.Longitude)
}

// UnmarshalJSON implements the json unmarshal to be used by the json-library
func (s *OutlookCoordinates) UnmarshalJSON(data []byte) error {
	tmp := struct {
		Accuracy    float64 `json:"accuracy"`
		Altitude    float64 `json:"altitude"`
		AltitudeAccuracy float64 `json:"altitudeAccuracy"`
		Latitude	float64`json:"latitude"`
		Longitude	float64 `json:"longitude"`
	}{}

	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}


	s.Longitude = tmp.Longitude
	s.Latitude = tmp.Latitude
	s.Accuracy = tmp.Accuracy
	s.Altitude = tmp.Altitude
	s.AltitudeAccuracy = tmp.AltitudeAccuracy

	return nil
}
