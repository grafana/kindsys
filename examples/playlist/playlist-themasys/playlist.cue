package playlist

name:        "Playlist"
maturity:    "merged"
description: "A playlist is a series of dashboards that is automatically rotated in the browser, on a configurable interval."

lineage: schemas: [
	{
	///////////////////////////////////////
	// Original internal ID only
	version: [0, 0]
	schema: {
		spec: {
			// Unique playlist identifier. Generated on creation, either by the
			// creator of the playlist of by the application.
			uid: string

			// Name of the playlist.
			name: string

			// Interval sets the time between switching views in a playlist.
			// FIXME: Is this based on a standardized format or what options are available? Can datemath be used?
			interval: string | *"5m"

			// The ordered list of items that the playlist will iterate over.
			// FIXME! This should not be optional, but changing it makes the godegen awkward
			items?: [...#PlaylistItem]
		} @cuetsy(kind="interface")

		#PlaylistItem: {
			// Type of the item.
			type: "dashboard_by_id" | "dashboard_by_tag"

			// Value depends on type and describes the playlist item.
			value: string

			// Title is used to display the item value
			title: string
		} @cuetsy(kind="interface")
	}
}, {
	//////////////////////////////////////////////////
	// Adding dashboard_by_uid, make title optional
	// deprecate dashboard_by_id
	version: [0, 1]
	schema: {
		spec: {
			// Unique playlist identifier. Generated on creation, either by the
			// creator of the playlist of by the application.
			// NOTE: uid MUST equal the k8s resource name
			uid: string

			// Name of the playlist.
			name: string

			// Interval sets the time between switching views in a playlist.
			// FIXME: Is this based on a standardized format or what options are available? Can datemath be used?
			interval: string | *"5m"

			// The ordered list of items that the playlist will iterate over.
			// FIXME! This should not be optional, but changing it makes the godegen awkward
			items?: [...#PlaylistItem]
		} @cuetsy(kind="interface")

		///////////////////////////////////////
		// Definitions (referenced above) are declared below

		#PlaylistItem: {
			// Type of the item.
			type: "dashboard_by_uid" | "dashboard_by_id" | "dashboard_by_tag"
			// Value depends on type and describes the playlist item.
			//
			//  - dashboard_by_id: The value is an internal numerical identifier set by Grafana. This
			//  is not portable as the numerical identifier is non-deterministic between different instances.
			//  Will be replaced by dashboard_by_uid in the future. (deprecated)
			//  - dashboard_by_tag: The value is a tag which is set on any number of dashboards. All
			//  dashboards behind the tag will be added to the playlist.
			//  - dashboard_by_uid: The value is the dashboard UID
			value: string

			// Title is an unused property -- it will be removed in the future
			title?: string
		} @cuetsy(kind="interface")
	}
},{
	//////////////////////////////////////////////////
	// Remove internal ID, remove title
	// and move the "uid" to the k8s "name" field
	version: [1, 0]
	schema: {
		spec: {
			// Name of the playlist.
			name: string

			// Interval sets the time between switching views in a playlist.
			// FIXME: Is this based on a standardized format or what options are available? Can datemath be used?
			interval: string | *"5m"

			// The ordered list of items that the playlist will iterate over.
			// FIXME! This should not be optional, but changing it makes the godegen awkward
			items?: [...#PlaylistItem]

			// Adding a required new field...
			// This is only hear so that thema breaking change detection allows
			// defining this as a new major version
			xxx: string
		} @cuetsy(kind="interface")

		///////////////////////////////////////
		// Definitions (referenced above) are declared below

		#PlaylistItem: {
			// Type of the item.
			type: "dashboard_by_uid" | "dashboard_by_id" | "dashboard_by_tag"
			// Value depends on type and describes the playlist item.
			//
			//  - dashboard_by_id: The value is an internal numerical identifier set by Grafana. This
			//  is not portable as the numerical identifier is non-deterministic between different instances.
			//  Will be replaced by dashboard_by_uid in the future. (deprecated)
			//  - dashboard_by_tag: The value is a tag which is set on any number of dashboards. All
			//  dashboards behind the tag will be added to the playlist.
			//  - dashboard_by_uid: The value is the dashboard UID
			value: string

			// Title is an unused property -- it will be removed in the future
			title?: string
		} @cuetsy(kind="interface")
	}
}]
