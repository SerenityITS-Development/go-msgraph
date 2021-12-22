package msgraph

import (
	"encoding/json"
	"fmt"
	"strings"
)

type OutlookCategories []OutlookCategory

func (t OutlookCategories) setGraphClient(u *User) OutlookCategories {
	for i := range t {
		t[i] = t[i].setGraphClient(u)
	}
	return t
}

func (t OutlookCategories) String() string {
	var categories = make([]string, len(t))
	for i, category := range t {
		categories[i] =  category.String()
	}
	return fmt.Sprintf("OutlookCategories(%v)", strings.Join(categories, ", "))
}

func (t OutlookCategories) FindCategoryByName(value string) (*OutlookCategory, error) {
	for _, category := range t {
		if category.DisplayName == value {
			return &category, nil
		}
	}
	return nil, ErrFindOutlookCategory
}

func (t *OutlookCategories) UnmarshalJSON(data []byte) error {
	tmp := struct {
		OutlookCategories []OutlookCategory `json:"value"`
	}{}
	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return fmt.Errorf("cannot UnmarshalJSON: %v | Data: %v", err, string(data))
	}
	*t = tmp.OutlookCategories
	return nil
}