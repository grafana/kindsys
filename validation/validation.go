package validation

import (
	"cuelang.org/go/cue"
	"fmt"
	"strings"
)

// EnsureNoExportedKindName checks that Kind and #Kind is prohibited as a field name in schemas
func EnsureNoExportedKindName(value cue.Value) error {
	sch := schema(value)

	var pos []string
	Walk(sch, func(v cue.Value) bool {
		label, _ := v.Label()
		label = strings.ToLower(label)
		if label == "#kind" || label == "kind" {
			pos = append(pos, v.Pos().String())
			return false
		}

		return true
	}, nil)

	if len(pos) != 0 {
		return fmt.Errorf("schema must not use `kind` keyword present at %s", strings.Join(pos, "; "))
	}

	return nil
}

func schema(v cue.Value) cue.Value {
	var sch cue.Value

	Walk(v, func(v cue.Value) bool {
		label, _ := v.Label()

		if label != "schemas" {
			return true
		}

		sch = v
		return false
	}, nil)

	return sch
}

// Copied from https://github.com/hofstadter-io/cuetils/

// Walk is an alternative to cue.Value.Walk which handles more field types
// You can customize this with your own options
func Walk(v cue.Value, before func(cue.Value) bool, after func(cue.Value), options ...cue.Option) {

	// call before and possibly stop recursion
	if before != nil && !before(v) {
		return
	}

	// possibly recurse
	switch v.IncompleteKind() {
	case cue.StructKind:
		if options == nil {
			options = defaultWalkOptions
		}
		s, _ := v.Fields(options...)

		for s.Next() {
			Walk(s.Value(), before, after, options...)
		}

	case cue.ListKind:
		l, _ := v.List()
		for l.Next() {
			Walk(l.Value(), before, after, options...)
		}
	}

	if after != nil {
		after(v)
	}
}

var defaultWalkOptions = []cue.Option{
	cue.Attributes(true),
	cue.Concrete(false),
	cue.Definitions(true),
	cue.Hidden(true),
	cue.Optional(true),
	cue.Docs(true),
}
