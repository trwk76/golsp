package main

import (
	"fmt"
	"sort"
	"strings"

	golang "github.com/trwk76/go-code/go"
)

func ModelFromData(data MetaModel) *Model {
	res := &Model{
		types: make(map[string]ModelType),
	}

	for _, itm := range data.Structures {
		res.types[itm.Name] = &StructType{
			mdl:  res,
			data: itm,
		}
	}

	for _, itm := range data.Enumerations {
		res.types[itm.Name] = &EnumType{
			mdl:  res,
			data: itm,
		}
	}

	for _, itm := range data.TypeAliases {
		res.types[itm.Name] = &AliasType{
			mdl:  res,
			data: itm,
		}
	}

	return res
}

type (
	Model struct {
		types map[string]ModelType
	}

	ModelType interface {
		golang() golang.TypeDecl
	}

	StructType struct {
		mdl  *Model
		data Structure
	}

	EnumType struct {
		mdl  *Model
		data Enumeration
	}

	AliasType struct {
		mdl  *Model
		data TypeAlias
	}
)

func (m *Model) goTypes() golang.TypeDecls {
	res := make(golang.TypeDecls, 0)

	for _, typ := range m.types {
		res = append(res, typ.golang())
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].ID < res[j].ID
	})

	return res
}

func (t *StructType) golang() golang.TypeDecl {
	bases := make([]golang.Type, len(t.data.Extends))

	for idx, ext := range t.data.Extends {
		if ref, ok := ext.Impl.(*ReferenceType); ok {
			bases[idx] = golang.Symbol{ID: golang.ID(ref.Name)}
		} else {
			panic(fmt.Errorf("struct %s: base #%d not a referrnce", t.data.Name, idx))
		}
	}

	return golang.TypeDecl{
		Comment: comment(t.data.Name, t.data.Documentation),
		ID:      golang.ID(t.data.Name),
		Spec:    golang.StructType{
			Bases: bases,
		},
	}
}

func (t *EnumType) golang() golang.TypeDecl {
	return golang.TypeDecl{
		Comment: comment(t.data.Name, t.data.Documentation),
		ID:      golang.ID(t.data.Name),
		Spec:    golang.String,
	}
}

func (t *AliasType) golang() golang.TypeDecl {
	return golang.TypeDecl{
		Comment: comment(t.data.Name, t.data.Documentation),
		ID:      golang.ID(t.data.Name),
		Spec:    golang.TypeAlias{Target: golang.String},
	}
}

func comment(name string, doc string) golang.Comment {
	buf := strings.Builder{}

	fmt.Fprintf(&buf, " %s\n", name)

	lines := strings.Split(strings.ReplaceAll(doc, "\r\n", "\n"), "\n")

	for idx, line := range lines {
		if len(line) > 0 {
			buf.WriteByte(' ')
			buf.WriteString(line)
		}

		if idx < len(lines) - 1 {
			buf.WriteByte('\n')
		}
	}

	return golang.Comment(buf.String())
}

var (
	_ ModelType = (*StructType)(nil)
	_ ModelType = (*EnumType)(nil)
	_ ModelType = (*AliasType)(nil)
)