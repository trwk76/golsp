package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func LoadModel() MetaModel {
	var mdl MetaModel

	raw, err := os.ReadFile("metaModel.json")
	if err != nil {
		panic(fmt.Errorf("error reading metaModel.json file: %s", err.Error()))
	}

	if err := json.Unmarshal(raw, &mdl); err != nil {
		panic(fmt.Errorf("error unmarshaling metaModel.json: %s", err.Error()))
	}

	return mdl
}

type (
	MetaModel struct {
		MetaData      MetaData       `json:"metaData"`
		Requests      []Request      `json:"requests"`
		Notifications []Notification `json:"notifications"`
		Structures    []Structure    `json:"structures"`
		Enumerations  []Enumeration  `json:"enumerations"`
		TypeAliases   []TypeAlias    `json:"typeAliases"`
	}

	MetaData struct {
		Version string `json:"version"`
	}

	Request struct {
		Deprecated          string           `json:"deprecated,omitempty"`
		Documentation       string           `json:"documentation,omitempty"`
		ErrorData           *Type            `json:"errorData,omitempty"`
		MessageDirection    MessageDirection `json:"messageDirection"`
		Method              string           `json:"method"`
		Params              OneOrMore[Type]  `json:"params,omitempty"`
		PartialResult       *Type            `json:"partialResult,omitempty"`
		Proposed            bool             `json:"proposed,omitempty"`
		RegistrationMethod  string           `json:"registrationMethod,omitempty"`
		RegistrationOptions *Type            `json:"registrationOptions,omitempty"`
		Result              Type             `json:"result"`
		Since               string           `json:"since,omitempty"`
	}

	Notification struct {
		Deprecated          string           `json:"deprecated,omitempty"`
		Documentation       string           `json:"documentation,omitempty"`
		MessageDirection    MessageDirection `json:"messageDirection"`
		Method              string           `json:"method"`
		Params              OneOrMore[Type]  `json:"params,omitempty"`
		Proposed            bool             `json:"proposed,omitempty"`
		RegistrationMethod  string           `json:"registrationMethod,omitempty"`
		RegistrationOptions *Type            `json:"registrationOptions,omitempty"`
		Since               string           `json:"since,omitempty"`
	}

	MessageDirection string

	Structure struct {
		Deprecated    string     `json:"deprecated,omitempty"`
		Documentation string     `json:"documentation,omitempty"`
		Extends       []Type     `json:"extends,omitempty"`
		Mixins        []Type     `json:"mixins,omitempty"`
		Name          string     `json:"name"`
		Properties    []Property `json:"properties"`
		Proposed      bool       `json:"proposed,omitempty"`
		Since         string     `json:"since,omitempty"`
	}

	Property struct {
		Deprecated    string `json:"deprecated,omitempty"`
		Documentation string `json:"documentation,omitempty"`
		Name          string `json:"name"`
		Optional      bool   `json:"optional,omitempty"`
		Proposed      bool   `json:"proposed,omitempty"`
		Since         string `json:"since,omitempty"`
		Type          Type   `json:"type"`
	}

	Enumeration struct {
		Deprecated           string             `json:"deprecated,omitempty"`
		Documentation        string             `json:"documentation,omitempty"`
		Name                 string             `json:"name"`
		Proposed             bool               `json:"proposed,omitempty"`
		Since                string             `json:"since,omitempty"`
		SupportsCustomValues bool               `json:"supportsCustomValues,omitempty"`
		Type                 Type               `json:"type"`
		Values               []EnumerationEntry `json:"values"`
	}

	EnumerationEntry struct {
		Deprecated    string           `json:"deprecated,omitempty"`
		Documentation string           `json:"documentation,omitempty"`
		Name          string           `json:"name"`
		Proposed      bool             `json:"proposed,omitempty"`
		Since         string           `json:"since,omitempty"`
		Value         EnumerationValue `json:"value"`
	}

	EnumerationValue struct {
		String string
		Number float64
	}

	TypeAlias struct {
		Deprecated    string `json:"deprecated,omitempty"`
		Documentation string `json:"documentation,omitempty"`
		Name          string `json:"name"`
		Proposed      bool   `json:"proposed,omitempty"`
		Since         string `json:"since,omitempty"`
		Type          Type   `json:"type"`
	}

	Type struct {
		Impl TypeImpl
	}

	TypeKind string

	TypeDisc struct {
		Kind TypeKind `json:"kind"`
	}

	TypeImpl interface {
	}

	AndType struct {
		TypeDisc
		Items []Type `json:"items"`
	}

	ArrayType struct {
		TypeDisc
		Element Type `json:"element"`
	}

	BaseType struct {
		TypeDisc
		Name string `json:"name"`
	}

	BooleanLiteralType struct {
		TypeDisc
		Value bool `json:"value"`
	}

	IntegerLiteralType struct {
		TypeDisc
		Value int64 `json:"value"`
	}

	MapType struct {
		TypeDisc
		Key   Type `json:"key"`
		Value Type `json:"value"`
	}

	OrType struct {
		TypeDisc
		Items []Type `json:"items"`
	}

	ReferenceType struct {
		TypeDisc
		Name string `json:"name"`
	}

	StringLiteralType struct {
		TypeDisc
		Value string `json:"value"`
	}

	StructLiteralType struct {
		TypeDisc
		Value StructLiteral `json:"value"`
	}

	TupleType struct {
		TypeDisc
		Items []Type `json:"items"`
	}

	StructLiteral struct {
		Deprecated    string     `json:"deprecated,omitempty"`
		Documentation string     `json:"documentation,omitempty"`
		Properties    []Property `json:"properties"`
		Proposed      bool       `json:"proposed,omitempty"`
		Since         string     `json:"since,omitempty"`
	}

	OneOrMore[T any] []T
)

const (
	ClientToServer MessageDirection = "clientToServer"
	ServerToClient MessageDirection = "serverToClient"
	Both           MessageDirection = "both"
)

const (
	TypeBase           TypeKind = "base"
	TypeReference      TypeKind = "reference"
	TypeArray          TypeKind = "array"
	TypeMap            TypeKind = "map"
	TypeAnd            TypeKind = "and"
	TypeOr             TypeKind = "or"
	TypeTuple          TypeKind = "tuple"
	TypeLiteral        TypeKind = "literal"
	TypeStringLiteral  TypeKind = "stringLiteral"
	TypeIntegerLiteral TypeKind = "integerLiteral"
	TypeBooleanLiteral TypeKind = "booleanLiteral"
)

func (v *EnumerationValue) UnmarshalJSON(raw []byte) error {
	if err := json.Unmarshal(raw, &v.Number); err == nil {
		return nil
	}

	return json.Unmarshal(raw, &v.String)
}

func (t *Type) UnmarshalJSON(raw []byte) error {
	var (
		disc TypeDisc
		impl TypeImpl
	)

	if err := json.Unmarshal(raw, &disc); err != nil {
		return err
	}

	impl = disc.newImpl()

	if err := json.Unmarshal(raw, &impl); err != nil {
		return err
	}

	t.Impl = impl
	return nil
}

func (d TypeDisc) newImpl() TypeImpl {
	switch d.Kind {
	case TypeBase:
		return &BaseType{}
	case TypeReference:
		return &ReferenceType{}
	case TypeArray:
		return &ArrayType{}
	case TypeMap:
		return &MapType{}
	case TypeAnd:
		return &AndType{}
	case TypeOr:
		return &OrType{}
	case TypeTuple:
		return &TupleType{}
	case TypeLiteral:
		return &StructLiteralType{}
	case TypeStringLiteral:
		return &StringLiteralType{}
	case TypeIntegerLiteral:
		return &IntegerLiteralType{}
	case TypeBooleanLiteral:
		return &BooleanLiteralType{}
	}

	panic(fmt.Errorf("type kind '%s' is not supported", d.Kind))
}

func (v *OneOrMore[T]) UnmarshalJSON(raw []byte) error {
	var (
		arr []T
		one T
	)

	if err := json.Unmarshal(raw, &arr); err == nil {
		*v = OneOrMore[T](arr)
		return nil
	}

	if err := json.Unmarshal(raw, &one); err != nil {
		return err
	}

	*v = OneOrMore[T]{one}
	return nil
}

var (
	_ TypeImpl         = (*AndType)(nil)
	_ TypeImpl         = (*ArrayType)(nil)
	_ TypeImpl         = (*BaseType)(nil)
	_ TypeImpl         = (*BooleanLiteralType)(nil)
	_ TypeImpl         = (*IntegerLiteralType)(nil)
	_ TypeImpl         = (*MapType)(nil)
	_ TypeImpl         = (*OrType)(nil)
	_ TypeImpl         = (*ReferenceType)(nil)
	_ TypeImpl         = (*StringLiteralType)(nil)
	_ TypeImpl         = (*StructLiteralType)(nil)
	_ TypeImpl         = (*TupleType)(nil)
	_ json.Unmarshaler = (*EnumerationValue)(nil)
	_ json.Unmarshaler = (*Type)(nil)
	_ json.Unmarshaler = (*OneOrMore[string])(nil)
)
