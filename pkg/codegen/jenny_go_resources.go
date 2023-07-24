package codegen

import (
	"bytes"
	"fmt"
	"go/format"
	"strings"

	"cuelang.org/go/cue"
	"github.com/dave/dst/dstutil"
	"github.com/grafana/codejen"
	"github.com/grafana/kindsys"
	"github.com/grafana/thema/encoding/gocode"
	"github.com/grafana/thema/encoding/openapi"
)

var schPath = cue.MakePath(cue.Hid("_#schema", "github.com/grafana/thema"))

type ResourceGoTypesJenny struct {
	ApplyFuncs       []dstutil.ApplyFunc
	ExpandReferences bool
}

func (*ResourceGoTypesJenny) JennyName() string {
	return "GoTypesJenny"
}

func (ag *ResourceGoTypesJenny) Generate(kind kindsys.Kind) (*codejen.File, error) {
	comm := kind.Props().Common()
	sfg := SchemaForGen{
		Name:    comm.Name,
		Schema:  kind.Lineage().Latest(),
		IsGroup: comm.LineageIsGroup,
	}
	sch := sfg.Schema

	iter, err := sch.Underlying().LookupPath(schPath).Fields()
	if err != nil {
		return nil, err
	}

	var subr []string
	for iter.Next() {
		subr = append(subr, typeNameFromKey(iter.Selector().String()))
	}

	buf := new(bytes.Buffer)
	mname := kind.Props().Common().MachineName
	if err := tmpls.Lookup("core_resource.tmpl").Execute(buf, tvars_resource{
		PackageName:      mname,
		KindName:         kind.Props().Common().Name,
		Version:          strings.Replace(sfg.Schema.Version().String(), ".", "-", -1),
		SubresourceNames: subr,
	}); err != nil {
		return nil, fmt.Errorf("failed executing core resource template: %w", err)
	}

	if err != nil {
		return nil, err
	}

	content, err := format.Source(buf.Bytes())
	if err != nil {
		return nil, err
	}

	return codejen.NewFile(fmt.Sprintf("pkg/kinds/%s/%s_gen.go", mname, mname), content, ag), nil
}

type SubresourceGoTypesJenny struct {
	ApplyFuncs       []dstutil.ApplyFunc
	ExpandReferences bool
}

func (*SubresourceGoTypesJenny) JennyName() string {
	return "GoResourceTypes"
}

func (g *SubresourceGoTypesJenny) Generate(kind kindsys.Kind) (codejen.Files, error) {
	comm := kind.Props().Common()
	sfg := SchemaForGen{
		Name:    comm.Name,
		Schema:  kind.Lineage().Latest(),
		IsGroup: comm.LineageIsGroup,
	}
	sch := sfg.Schema

	// Iterate through all top-level fields and make go types for them
	// (this should consist of "spec" and arbitrary subresources)
	i, err := sch.Underlying().LookupPath(schPath).Fields()
	if err != nil {
		return nil, err
	}
	files := make(codejen.Files, 0)
	for i.Next() {
		str := i.Selector().String()

		b, err := gocode.GenerateTypesOpenAPI(sch, &gocode.TypeConfigOpenAPI{
			// TODO will need to account for sanitizing e.g. dashes here at some point
			Config: &openapi.Config{
				Group:    false, // TODO: better
				RootName: typeNameFromKey(str),
				Subpath:  cue.MakePath(cue.Str(str)),
			},
			PackageName: sfg.Schema.Lineage().Name(),
			ApplyFuncs:  append(g.ApplyFuncs, PrefixDropper(sfg.Name)),
		})
		if err != nil {
			return nil, err
		}

		name := sfg.Schema.Lineage().Name()
		files = append(files, codejen.File{
			RelativePath: fmt.Sprintf("pkg/kinds/%s/%s_%s_gen.go", name, name, strings.ToLower(str)),
			Data:         b,
			From:         []codejen.NamedJenny{g},
		})
	}

	return files, nil
}

func typeNameFromKey(key string) string {
	if len(key) > 0 {
		return strings.ToUpper(key[:1]) + key[1:]
	}
	return strings.ToUpper(key)
}
