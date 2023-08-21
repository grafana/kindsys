package v1x

import "fmt"

//*******************************************************************************************
// NOTE!!
// This file is exploring generating the JSONSchema from golang... but that is paused for now
// The tests are just about the schema for now
//*******************************************************************************************

// Type of the item.
type ItemType string

// Defines values for ItemType.
const (
	ItemTypeDashboardByTag ItemType = "dashboard_by_tag"
	ItemTypeDashboardByUid ItemType = "dashboard_by_uid"
)

// Item defines model for Item.
type Item struct {
	// Type of the item.
	Type ItemType `json:"type"`

	// Value depends on type and describes the playlist item.
	//
	//  - dashboard_by_id: The value is an internal numerical identifier set by Grafana. This
	//  is not portable as the numerical identifier is non-deterministic between different instances.
	//  Will be replaced by dashboard_by_uid in the future. (deprecated)
	//  - dashboard_by_tag: The value is a tag which is set on any number of dashboards. All
	//  dashboards behind the tag will be added to the playlist.
	//  - dashboard_by_uid: The value is the dashboard UID
	Value string `json:"value"`
}

// Spec defines model for Spec.
type Spec struct {
	// Interval sets the time between switching views in a playlist.
	// FIXME: Is this based on a standardized format or what options are available? Can datemath be used?
	Interval string `json:"interval"`

	// The ordered list of items that the playlist will iterate over.
	Items []Item `json:"items"`

	// Name of the playlist.
	Name string `json:"name"`

	// XXX is just used so thema can detect a breaking change at version 1.0
	Xxx string `json:"xxx,omitempty"`
}

func (s *Spec) Validate(major int, minor int) error {
	if major != 1 {
		return fmt.Errorf("expecting major version = 1")
	}
	if len(s.Name) < 1 {
		return fmt.Errorf("missing name")
	}
	if len(s.Items) < 1 {
		return fmt.Errorf("missing items")
	}
	for idx, item := range s.Items {
		if len(item.Value) < 1 {
			return fmt.Errorf("missing item value (item[%d])", idx)
		}
		switch item.Type {
		case ItemTypeDashboardByTag:
		case ItemTypeDashboardByUid:
		default:
			return fmt.Errorf("invalid item type (item[%d])", idx)
		}
	}
	return nil
}
