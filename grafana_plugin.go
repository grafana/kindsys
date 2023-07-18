package kindsys

import (
	"cuelang.org/go/cue"
	"cuelang.org/go/cue/errors"
	"fmt"
	"github.com/grafana/thema"
)

var allSi map[string]SchemaInterface

func init() {
	allSi = make(map[string]SchemaInterface, 0)
	for n, s := range SchemaInterfaces(nil) {
		allSi[n] = s
	}
}

func ToGrafanaPluginComposable(rt *thema.Runtime, v cue.Value, skipSchema bool) ([]Composable, error) {
	fw := CUEFramework(rt.Context())
	gp := fw.LookupPath(cue.ParsePath("GrafanaPlugin"))
	if !gp.Exists() {
		return nil, errors.New("Cannot find GrafanaPlugin template")
	}

	unify := v
	if !skipSchema {
		unify = v.Unify(gp)
		if unify.Err() != nil {
			return nil, errors.Newf(v.Pos(), "Schema doesn't follow GrafanaPlugin pattern: %s", unify.Err())
		}
	}

	compo := unify.LookupPath(cue.ParsePath("composableKinds"))
	if !compo.Exists() {
		return nil, errors.New("Composable plugin should include `composableKinds` root struct")
	}

	compoList := make([]Composable, 0)
	iter, _ := compo.Fields()
	for iter.Next() {
		contract, ok := allSi[iter.Selector().String()]
		if !ok {
			return nil, errors.Newf(iter.Value().Pos(), "Unable to find %s schema interface", iter.Selector().String())
		}

		f := compo.LookupPath(cue.MakePath(iter.Selector()))
		schemas := f.LookupPath(cue.MakePath(cue.Str("lineage"), cue.Str("schemas")))
		fmt.Println(schemas)
		schList, _ := schemas.List()
		for schList.Next() {
			sch := schList.Value().LookupPath(cue.ParsePath("schema"))
			if err := sch.Subsume(contract.Contract()); err != nil {
				return nil, err
			}
		}

		propsDef, err := ToKindProps[ComposableProperties](f)
		if err != nil {
			return nil, errors.Newf(f.Pos(), "Cannot parse properties: %s", err)
		}

		c, err := BindComposable(rt, propsDef)
		if err != nil {
			return nil, err
		}

		compoList = append(compoList, c)
	}

	return compoList, nil
}
