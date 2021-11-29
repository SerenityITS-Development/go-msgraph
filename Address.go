package msgraph

import (
	"encoding/json"
	"fmt"
)

type Address struct {
	City	string `json:"city,omitempty"`
	CountryOrRegion	string `json:"countryOrRegion,omitempty"`
	PostalCode	string `json:"postalCode,omitempty"`
	State    string `json:"state,omitempty"`
	Street   string `json:"street,omitempty"`
}

func (s Address) String() string {
	return fmt.Sprintf("Address(Street: %s, City: %s, State: %s, PostalCode: %s, CountryOrRegion: %s)",
		s.Street, s.City, s.State, s.PostalCode, s.CountryOrRegion)
}

// Equal compares the TimeSlot to the other value and returns true
// if the values are equal
func (s Address) Equal(other Address) bool {
	return s.City == other.City && s.CountryOrRegion == other.CountryOrRegion &&
		s.State == other.State && s.Street == other.Street && s.PostalCode == other.PostalCode
}

func (s *Address) UnmarshalJSON(data []byte) error {
	tmp := struct {
		City	string `json:"city"`
		CountryOrRegion	string `json:"countryOrRegion"`
		PostalCode	string `json:"postalCode"`
		State    string `json:"state"`
		Street   string `json:"street"`
	}{}

	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}

	s.PostalCode = tmp.PostalCode
	s.City = tmp.City
	s.Street = tmp.Street
	s.State = tmp.State
	s.CountryOrRegion = tmp.CountryOrRegion

	return nil
}