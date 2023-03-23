package validation

import (
	"cuelang.org/go/cue"
	"cuelang.org/go/cue/ast"
	"fmt"
	"strings"
)

// EnsureNoExportedKindName checks that Kind is prohibited as a name in the following cases:
// #Kind: _
// Kind: _ @cuetsy(*)
func EnsureNoExportedKindName(value cue.Value) error {
	lin := lineage(value.Source())

	position, found := searchKindKeyword(lin)
	if found {
		return fmt.Errorf("schema must not use `Kind` keyword present at %s", position)
	}

	return nil
}

func searchKindKeyword(lin ast.Node) (string, bool) {
	var pos string
	var found bool

	ast.Walk(lin, func(n ast.Node) bool {
		field, is := n.(*ast.Field)
		if !is {
			return true
		}

		label, is := field.Label.(*ast.Ident)
		if !is {
			return true
		}

		var isAttr bool
		for _, a := range field.Attrs {
			if strings.Contains(a.Text, "@cuetsy") {
				isAttr = true
			}
		}
		if label.String() == "#Kind" || (label.String() == "Kind" && isAttr) {
			found = true
			pos = label.Pos().Position().String()
			return false
		}
		return true
	}, nil)

	return pos, found
}

func lineage(node ast.Node) ast.Node {
	var lin ast.Node

	ast.Walk(node, func(n ast.Node) bool {
		field, is := n.(*ast.Field)
		if !is {
			return true
		}

		label, is := field.Label.(*ast.Ident)
		if !is {
			return true
		}

		if label.String() != "lineage" {
			return true
		}

		lin = n
		return false
	}, nil)

	return lin
}
