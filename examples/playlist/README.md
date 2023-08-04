# Playlist examples
This package explores how kindsys can implement the playlist kind.

### Why playlist?  

It is a silly, but simple and illustrative example that can exercise real world requirements.

### History

Playlists have been in core grafana for a long time. The APIs were originally written against 
internal dashboard ids (numeric integers).  
* In v8.? (TODO) -- we allowed saving references to UIDs in addition to internal IDs
* In v9.? (TODO) -- we updated the UI so that the "title" attribute is not used, it now shows the loaded dashboard title instead.  This should also ensure that the spec.uid field is the k8s wrapper name  
* In some future version -- we want to force everythign do reference dashboard UIDs and remove the unused title version.  This version will also remove the uid field on spec because that is a duplicate for what is in the wrapper metadata


## Schemas

In this example, we will define three versions

* v0.0 -- the original schema that only takes internal ids
* v0.1 -- adds a uid option to each playlist item and deprecates the id version and unused title
* v1.0 -- removes the internal ID option, and unused title properties

### Compatibility

#### v0.0 -> v0.1
✅ This just adds additional features that may or may not be used

#### v0.1 -> v1.0
⚠️ Requires database access to convert ID to UID

#### v1.0 -> v0.1
✅ This can migrate OK

#### v0.1 -> v0.0
⚠️ Requires database access to convert UID to ID
⚠️ Requires database access to lookup title from UID





