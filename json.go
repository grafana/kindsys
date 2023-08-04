package kindsys

import (
	"fmt"
	"io"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
)

// weird that we may want this public... but ¯\_(ツ)_/¯ lets work with it for now
type JSONResourceBuilder struct {
	// Called when the group+version+kind have been read
	SetGroupVersionKind func(group, version, kind string) error

	// Called when a non-standard annotation field is found
	// (before the metadata callback)
	SetAnnotation func(key string, val string)

	// Called after the metadata tag is closed
	SetMetadata func(s StaticMetadata, c CommonMetadata)

	// Called when the parse finds a spec element
	ReadSpec   func(iter *jsoniter.Iterator) error
	ReadStatus func(iter *jsoniter.Iterator) error
	ReadSub    func(name string, iter *jsoniter.Iterator) error
}

func ReadResourceJSON(reader io.Reader, builder JSONResourceBuilder) error {
	iter := jsoniter.Parse(jsoniter.ConfigDefault, reader, 1024)

	static := StaticMetadata{}

	for l1Field := iter.ReadObject(); l1Field != ""; l1Field = iter.ReadObject() {
		err := iter.Error
		switch l1Field {
		case "apiVersion":
			err = static.setGroupVersionFromAPI(iter.ReadString())

		case "kind":
			static.Kind = iter.ReadString()
			if err == nil && builder.SetGroupVersionKind != nil {
				err = builder.SetGroupVersionKind(static.Group, static.Version, static.Kind)
			}

		case "metadata":
			common := CommonMetadata{}
			for l2Field := iter.ReadObject(); l2Field != ""; l2Field = iter.ReadObject() {
				switch l2Field {
				case "namespace":
					static.Namespace = iter.ReadString()
				case "name":
					static.Name = iter.ReadString()
				case "resourceVersion":
					common.ResourceVersion = iter.ReadString()
				case "uid":
					common.UID = iter.ReadString()
				case "creationTimestamp":
					t, err := time.Parse(time.RFC3339, iter.ReadString())
					if err != nil {
						return fmt.Errorf("invalid updatedTimestamp format // %w", err)
					}
					common.CreationTimestamp = t
				case "finalizers":
					common.Finalizers = make([]string, 0)
					for iter.ReadArray() {
						common.Finalizers = append(common.Finalizers, iter.ReadString())
					}

				case "annotations":
					for anno := iter.ReadObject(); anno != ""; anno = iter.ReadObject() {
						val := iter.ReadString()

						i := strings.Index(anno, "/")
						if i > 0 {
							g := anno[:i]
							v := anno[i+1:]

							if g == "grafana.com" {
								switch v {
								case "createdBy":
									common.CreatedBy = val
								case "updatedBy":
									common.UpdatedBy = val
								case "creationTimestamp":
									t, err := time.Parse(time.RFC3339, val)
									if err != nil {
										return fmt.Errorf("invalid creationTimestamp format // %w", err)
									}
									common.CreationTimestamp = t
								case "updatedTimestamp":
									t, err := time.Parse(time.RFC3339, val)
									if err != nil {
										return fmt.Errorf("invalid updatedTimestamp format // %w", err)
									}
									common.UpdateTimestamp = t
								case "origin.name":
									if common.Origin == nil {
										common.Origin = &ResourceOriginInfo{}
									}
									common.Origin.Name = val
								case "origin.path":
									if common.Origin == nil {
										common.Origin = &ResourceOriginInfo{}
									}
									common.Origin.Path = val
								case "origin.key":
									if common.Origin == nil {
										common.Origin = &ResourceOriginInfo{}
									}
									common.Origin.Key = val
								case "origin.timestamp":
									if common.Origin == nil {
										common.Origin = &ResourceOriginInfo{}
									}
									t, err := time.Parse(time.RFC3339, val)
									if err != nil {
										return fmt.Errorf("invalid updatedTimestamp format // %w", err)
									}
									common.Origin.Timestamp = &t
								default:
									fmt.Printf("grafana anno> %s = %v\n", g, v)
								}
							} else {
								fmt.Printf("anno ???> %s = %v\n", anno, val)
							}

						} else {
							switch anno {
							default:
								fmt.Printf("anno> %s = %v\n", anno, val)
							}
						}
						if iter.Error != nil {
							return iter.Error
						}
					}
				default:
					tt := iter.ReadAny()
					fmt.Printf("meta> %s = %v\n", l2Field, tt)
					if common.ExtraFields == nil {
						common.ExtraFields = make(map[string]any)
					}
					common.ExtraFields[l2Field] = tt
				}
				if iter.Error != nil {
					return iter.Error
				}
			}
			builder.SetMetadata(static, common)

		case "spec":
			err = builder.ReadSpec(iter)
		case "status":
			if builder.ReadStatus == nil {
				return fmt.Errorf("unsupported subresource")
			}
			err = builder.ReadStatus(iter)
		default:
			if builder.ReadSub == nil {
				return fmt.Errorf("unsupported subresource")
			}
			err = builder.ReadSub(l1Field, iter)
		}
		if err != nil {
			return err
		}
		if iter.Error != nil {
			return iter.Error
		}
	}
	return iter.Error
}

func WriteResourceJSON(obj Resource, stream *jsoniter.Stream) error {
	isMore := false
	static := obj.StaticMetadata()
	common := obj.CommonMetadata()
	custom := obj.CustomMetadata() // ends up in annotations

	stream.WriteObjectStart()
	stream.WriteObjectField("apiVersion")
	stream.WriteString(static.GetAPIVersion())
	stream.WriteMore()
	stream.WriteObjectField("kind")
	stream.WriteString(static.Kind)

	stream.WriteMore()
	stream.WriteObjectField("metadata")
	stream.WriteObjectStart()
	isMore = writeOptionalString(false, "name", static.Name, stream)
	isMore = writeOptionalString(isMore, "namespace", static.Namespace, stream)

	if isMore {
		stream.WriteMore()
	}
	stream.WriteObjectField("annotations")
	stream.WriteObjectStart()

	prefix := "grafana.com/"
	isMore = writeOptionalString(false, prefix+"createdBy", common.CreatedBy, stream)
	isMore = writeOptionalString(isMore, prefix+"updatedBy", common.UpdatedBy, stream)

	origin := common.Origin
	if origin != nil && origin.Key != "" {
		isMore = writeOptionalString(isMore, prefix+"origin.name", origin.Name, stream)
		isMore = writeOptionalString(isMore, prefix+"origin.path", origin.Path, stream)
		isMore = writeOptionalString(isMore, prefix+"origin.key", origin.Key, stream)
		isMore = writeOptionalTime(isMore, prefix+"origin.timestamp", origin.Timestamp, stream)
	}
	if stream.Error != nil {
		return stream.Error
	}

	// Currently added directly to the annotations
	if custom != nil {
		for k, v := range custom.MapFields() {
			if isMore {
				stream.WriteMore()
			}
			stream.WriteObjectField(k)
			stream.WriteVal(v)
			isMore = true
		}
		if stream.Error != nil {
			return stream.Error
		}
	}

	_ = writeOptionalTime(isMore, prefix+"updatedTimestamp", &common.UpdateTimestamp, stream)
	stream.WriteObjectEnd() // annotations

	isMore = writeOptionalTime(true, "creationTimestamp", &common.CreationTimestamp, stream)
	isMore = writeOptionalTime(isMore, "deletionTimestamp", common.DeletionTimestamp, stream)
	isMore = writeOptionalString(isMore, "resourceVersion", common.ResourceVersion, stream)
	if len(common.Labels) > 0 {
		if isMore {
			stream.WriteMore()
		}
		stream.WriteObjectField("labels")
		stream.WriteVal(common.Labels)
		isMore = true
	}
	if len(common.Finalizers) > 0 {
		if isMore {
			stream.WriteMore()
		}
		stream.WriteObjectField("finalizers")
		stream.WriteArrayStart()
		for i, v := range common.Finalizers {
			if i > 0 {
				stream.WriteMore()
			}
			stream.WriteVal(v)
		}
		stream.WriteArrayEnd()
		isMore = true
	}
	isMore = writeOptionalString(isMore, "uid", common.UID, stream)

	// Additional k8s fields
	if common.ExtraFields != nil {
		for k, v := range common.ExtraFields {
			if isMore {
				stream.WriteMore()
			}
			stream.WriteObjectField(k)
			stream.WriteVal(v)
			isMore = true
		}
	}

	// End the metadata section
	stream.WriteObjectEnd()

	// SPEC
	spec := obj.SpecObject()
	if spec != nil {
		stream.WriteMore()
		stream.WriteObjectField("spec")
		stream.WriteVal(spec)
		if stream.Error != nil {
			return stream.Error
		}
	}

	// This will get status etc
	for k, v := range obj.Subresources() {
		if v != nil {
			stream.WriteMore()
			stream.WriteObjectField(k)
			stream.WriteVal(v)
		}
	}

	stream.WriteObjectEnd()
	// fmt.Printf("XXX:\n\n\n%s\n\n\n", string(stream.Buffer()))
	return stream.Error
}

func writeOptionalString(isMore bool, key string, val string, stream *jsoniter.Stream) bool {
	if val == "" {
		return isMore
	}
	if isMore {
		stream.WriteMore()
	}
	stream.WriteObjectField(key)
	stream.WriteString(val)
	return true
}

func writeOptionalTime(isMore bool, key string, val *time.Time, stream *jsoniter.Stream) bool {
	if val == nil {
		return isMore
	}
	if isMore {
		stream.WriteMore()
	}
	stream.WriteObjectField(key)
	stream.WriteString(val.Format(time.RFC3339))
	return true
}
