package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/trwk76/go-code"
	g "github.com/trwk76/go-code/go"
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
		golang(extra *g.Decls) g.TypeDecl
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

func (m *Model) decls() g.Decls {
	types := make(g.TypeDecls, 0)
	extra := make(g.Decls, 0)

	for _, typ := range m.types {
		types = append(types, typ.golang(&extra))
	}

	sort.Slice(types, func(i, j int) bool {
		return types[i].ID < types[j].ID
	})

	return append(
		g.Decls{types},
		extra...,
	)
}

func (m *Model) goType(t Type, opt bool) g.Type {

}

func (t *StructType) golang(extra *g.Decls) g.TypeDecl {
	bases := make([]g.Type, len(t.data.Extends))
	for idx, ext := range t.data.Extends {
		if ref, ok := ext.Impl.(*ReferenceType); ok {
			bases[idx] = g.Symbol{ID: g.ID(ref.Name)}
		} else {
			panic(fmt.Errorf("struct %s: base #%d not a referrnce", t.data.Name, idx))
		}
	}

	flds := make([]g.StructField, 0)

	for idx, ext := range t.data.Mixins {
		if ref, ok := ext.Impl.(*ReferenceType); ok {
			flds = append(flds, t.mdl.types[ref.Name].(*StructType).AllFields()...)
		} else {
			panic(fmt.Errorf("struct %s: mixin #%d not a referrnce", t.data.Name, idx))
		}
	}

	return g.TypeDecl{
		Comment: comment(t.data.Name, t.data.Documentation),
		ID:      g.ID(t.data.Name),
		Spec: g.StructType{
			Bases:  bases,
			Fields: append(flds, t.Fields()...),
		},
	}
}

func (t *StructType) Fields() []g.StructField {
	res := make([]g.StructField, len(t.data.Properties))

	for idx, prop := range t.data.Properties {
		res[idx] = g.StructField{
			ID:   g.ID(code.IDToPascal(prop.Name)),
			Type: t.goType(prop.Type, prop.Optional),
			Tags: g.Tags{{Name: "json", Value: prop.Name}},
		}
	}

	return res
}

func (t *StructType) AllFields() []g.StructField {
	res := make([]g.StructField, 0)

	for idx, ext := range t.data.Extends {
		if ref, ok := ext.Impl.(*ReferenceType); ok {
			res = append(res, t.mdl.types[ref.Name].(*StructType).AllFields()...)
		} else {
			panic(fmt.Errorf("struct %s: base #%d not a referrnce", t.data.Name, idx))
		}
	}

	for idx, ext := range t.data.Mixins {
		if ref, ok := ext.Impl.(*ReferenceType); ok {
			res = append(res, t.mdl.types[ref.Name].(*StructType).AllFields()...)
		} else {
			panic(fmt.Errorf("struct %s: mixin #%d not a referrnce", t.data.Name, idx))
		}
	}

	return append(res, t.Fields()...)
}

func (t *EnumType) golang(extra *g.Decls) g.TypeDecl {
	var (
		base  g.Type
		toVal func(val EnumerationValue) g.Expr
	)

	if btyp, ok := t.data.Type.Impl.(*BaseType); ok {
		switch btyp.Name {
		case "integer":
			base = g.Int32
			toVal = func(val EnumerationValue) g.Expr { return g.IntExpr(int32(val.Number)) }
		case "number":
			base = g.Float64
			toVal = func(val EnumerationValue) g.Expr { return g.FloatExpr(val.Number) }
		case "string":
			base = g.String
			toVal = func(val EnumerationValue) g.Expr { return g.StringExpr(val.String) }
		case "uinteger":
			base = g.Uint32
			toVal = func(val EnumerationValue) g.Expr { return g.UintExpr(uint32(val.Number)) }
		default:
			panic(fmt.Errorf("enumeration '%s': unsupported base type '%s'", t.data.Name, btyp.Name))
		}
	} else {
		panic(fmt.Errorf("enumeration '%s': unsupported base type", t.data.Name))
	}

	vals := make([]g.ConstDecl, len(t.data.Values))

	for idx, itm := range t.data.Values {
		vals[idx] = g.ConstDecl{
			Comment: g.Comment(itm.Documentation),
			ID:      g.ID(t.data.Name + itm.Name),
			Type:    g.Symbol{ID: g.ID(t.data.Name)},
			Value:   toVal(itm.Value),
		}
	}

	*extra = append(*extra, g.ConstDecls(vals))

	return g.TypeDecl{
		Comment: comment(t.data.Name, t.data.Documentation),
		ID:      g.ID(t.data.Name),
		Spec:    base,
	}
}

func (t *AliasType) golang(extra *g.Decls) g.TypeDecl {
	return g.TypeDecl{
		Comment: comment(t.data.Name, t.data.Documentation),
		ID:      g.ID(t.data.Name),
		Spec:    g.TypeAlias{Target: g.String},
	}
}

func comment(name string, doc string) g.Comment {
	buf := strings.Builder{}

	fmt.Fprintf(&buf, " %s\n", name)

	lines := strings.Split(strings.ReplaceAll(doc, "\r\n", "\n"), "\n")

	for idx, line := range lines {
		if len(line) > 0 {
			buf.WriteByte(' ')
			buf.WriteString(line)
		}

		if idx < len(lines)-1 {
			buf.WriteByte('\n')
		}
	}

	return g.Comment(buf.String())
}

var (
	_ ModelType = (*StructType)(nil)
	_ ModelType = (*EnumType)(nil)
	_ ModelType = (*AliasType)(nil)
)
