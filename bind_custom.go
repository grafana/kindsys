package kindsys

import (
	"fmt"
	"sort"

	"github.com/grafana/thema"
)

// genericCore is a dynamically typed representation of a parsed and
// validated [Custom] kind, implemented with thema.
type genericCustom struct {
	def Def[CustomProperties]
	lin thema.Lineage

	// map of string name of slot to the currently composed contents of the slot
	composed map[string][]Composable
}

func (k genericCustom) FromBytes(b []byte, codec Decoder) (*UnstructuredResource, error) {
	inst, err := bytesToAnyInstance(k, b, codec)
	if err != nil {
		return nil, err
	}
	// we have a valid instance! decode into unstructured
	// TODO implement me
	_ = inst
	panic("implement me")
}

func (k genericCustom) Validate(b []byte, codec Decoder) error {
	// TODO implement me
	panic("implement me")
}

func (k genericCustom) CurrentVersion() thema.SyntacticVersion {
	return k.def.Properties.CurrentVersion
}

func (k genericCustom) Group() string {
	return k.def.Properties.CRD.Group
}

var _ Custom = genericCustom{}

// Props returns the generic SomeKindProperties
func (k genericCustom) Props() SomeKindProperties {
	return k.def.Properties
}

// Name returns the Name property
func (k genericCustom) Name() string {
	return k.def.Properties.Name
}

// MachineName returns the MachineName property
func (k genericCustom) MachineName() string {
	return k.def.Properties.MachineName
}

// Maturity returns the Maturity property
func (k genericCustom) Maturity() Maturity {
	return k.def.Properties.Maturity
}

// Def returns a Def with the type of ExtendedProperties, containing the bound ExtendedProperties
func (k genericCustom) Def() Def[CustomProperties] {
	return k.def
}

// Lineage returns the underlying bound Lineage
func (k genericCustom) Lineage() thema.Lineage {
	return k.lin
}

func (k genericCustom) Compose(slot Slot, kinds ...Composable) (Custom, error) {
	// first, check that this kind supports this slot
	if k.def.Properties.Slots[slot.Name] != slot {
		return nil, &ErrNoSlotInKind{
			Slot: slot,
			Kind: k,
		}
	}

	schif, err := FindSchemaInterface(slot.SchemaInterface)
	if err != nil {
		panic(fmt.Sprintf("unreachable - slot was for nonexistent schema interface %s which should have been rejected at bind time", slot.SchemaInterface))
	}

	// then check that all provided kinds are implementors of the slot
	for _, kind := range kinds {
		if kind.Implements().Name() != schif.Name() {
			return nil, &ErrKindDoesNotImplementInterface{
				Kind:      kind,
				Interface: schif,
			}
		}
	}

	// Inputs look good. Make a copy with our built-up compose map
	com := make(map[string][]Composable)
	for k, v := range k.composed {
		com[k] = v
	}

	var all []Composable
	copy(all, com[slot.Name])
	all = append(all, kinds...)
	// Sort to ensure deterministic output of validation error messages, etc.
	sort.Slice(all, func(i, j int) bool {
		return all[i].Name() < all[j].Name()
	})
	com[slot.Name] = all

	return genericCustom{
		def:      k.def,
		lin:      k.lin,
		composed: com,
	}, nil
}

// BindCustom creates a Custom-implementing type from a def, runtime, and opts
//
//nolint:lll
func BindCustom(rt *thema.Runtime, def Def[CustomProperties], opts ...thema.BindOption) (Custom, error) {
	lin, err := def.Some().BindKindLineage(rt, opts...)
	if err != nil {
		return nil, err
	}

	return genericCustom{
		def: def,
		lin: lin,
	}, nil
}

// TODO docs
func BindCustomResource[R Resource](k Custom) (TypedCustom[R], error) {
	// TODO implement me
	panic("implement me")
}
