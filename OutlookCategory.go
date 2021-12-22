package msgraph

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type OutlookCategory struct {
	graphClient	*GraphClient
	user 		*User

	ID 			string 	`json:"id"`
	DisplayName string 	`json:"displayName"`
	Color 		string 	`json:"color"`
}


func (t OutlookCategory) setGraphClient(u *User) OutlookCategory {
	t.graphClient = u.graphClient
	t.user = u
	return t
}

func (t OutlookCategory) String() string {
	return fmt.Sprintf("OutlookCategory(ID: \"%v\", DisplayName: \"%v\", Color: \"%v\"", t.ID, t.DisplayName, t.Color)
}

func (t OutlookCategory) Delete(opts ...DeleteQueryOption) error {
	if t.graphClient == nil {
		return ErrNotGraphClientSourced
	}

	resource := fmt.Sprintf("/users/%v/outlook/masterCategories/%v", t.user.ID, t.ID)

	err := t.graphClient.makeDELETEAPICall(resource, compileDeleteQueryOptions(opts), nil)
	return err
}

func (t *OutlookCategory) Update(opts ...UpdateQueryOption) error {
	if t.graphClient == nil {
		return ErrNotGraphClientSourced
	}

	resource := fmt.Sprintf("/users/%v/outlook/masterCategories/%v", t.user.ID, t.ID)

	bodyBytes, err := json.Marshal(t)
	if err != nil {
		return err
	}
	reader := bytes.NewReader(bodyBytes)

	// TODO: check return body, maybe there is some potential success or error message hidden in it?
	err = t.graphClient.makePATCHAPICall(resource, compileUpdateQueryOptions(opts), reader, nil)
	return err
}