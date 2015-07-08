package parser

import (
	"bytes"
	"fmt"
	"io"
	"sort"
	"strings"
)

// Struct structure contains some properties that a GO struct need
type Struct struct {
	Name       string            // name of struct
	Properties []*StructProperty // list of property that belong to struct

	nameLen, typeLen int // max len of name and type
}

// NewStruct create new struct with given name.
// Name of a struct is mandatory.
func NewStruct(name string) *Struct {
	// todo: validate name
	s := &Struct{
		Name: name,
	}

	return s
}

// AddProperty append a struct property to list of properties of struct
func (s *Struct) AddProperty(sp ...*StructProperty) {
	if s == nil || sp == nil {
		return
	}

	for i := range sp {
		if len(sp[i].Name) > s.nameLen {
			s.nameLen = len(sp[i].Name)
		}

		if len(sp[i].Type) > s.typeLen {
			s.typeLen = len(sp[i].Type)
		}
	}

	s.Properties = append(s.Properties, sp...)
}

// WriteTo string of struct code to given IO
func (s *Struct) WriteTo(w io.Writer) (int64, error) {
	indent := []byte("\t")
	newline := []byte("\n")
	buffer := &bytes.Buffer{}

	buffer.WriteString("type " + s.Name + " struct {")
	buffer.Write(newline)

	// sort property by name
	sort.Sort(byName{s.Properties})

	for i := range s.Properties {
		buffer.Write(indent)
		buffer.WriteString(s.Properties[i].String(s.nameLen, s.typeLen))
		buffer.Write(newline)
	}

	buffer.WriteString("}")

	return buffer.WriteTo(w)
}

type structProperties []*StructProperty

// byName implements sort.Interface for []*StructProperty based on the Name field.
type byName struct{ structProperties }

func (n byName) Len() int {
	return len(n.structProperties)
}

func (n byName) Swap(i, j int) {
	n.structProperties[i], n.structProperties[j] = n.structProperties[j], n.structProperties[i]
}

func (n byName) Less(i, j int) bool {
	return n.structProperties[i].Name < n.structProperties[j].Name
}

// StructProperty define a property of struct
type StructProperty struct {
	Name string       // property name
	Type string       // type of property (todo: maybe a kind of reflect type)
	Tags PropertyTags // contains tag of a property
}

// NewStructProperty create a new struct property instance
func NewStructProperty(name, ptype string, tags PropertyTags) *StructProperty {
	sp := &StructProperty{
		Name: name,
		Type: ptype,
		Tags: tags,
	}

	return sp
}

// AddTag add a tag to property, return error if tag exists
func (sp *StructProperty) AddTag(name, value string) error {
	if sp.Tags == nil {
		sp.Tags = PropertyTags{}
	}

	if _, found := sp.Tags[name]; !found {
		sp.Tags[name] = value
		return nil
	}

	return fmt.Errorf("tag %q exists", name)
}

// String return a string of code for a property
func (sp *StructProperty) String(nameLen, typeLen int) string {
	if sp == nil {
		return ""
	}

	buf := &bytes.Buffer{}

	buf.WriteString(sp.Name)
	if nameLen > len(sp.Name) {
		buf.WriteString(strings.Repeat(" ", nameLen-len(sp.Name)+1))
	} else {
		buf.WriteString(" ")
	}

	buf.WriteString(sp.Type)
	if typeLen > len(sp.Type) {
		buf.WriteString(strings.Repeat(" ", typeLen-len(sp.Type)+1))
	} else {
		buf.WriteString(" ")
	}

	buf.WriteString("`" + sp.Tags.String() + "`")

	return buf.String()
}

// PropertyTags is a type of tag of a property
// maybe it can be map[string][]string?
type PropertyTags map[string]string

// String will return property tags into a string
func (pt PropertyTags) String() string {
	if pt == nil {
		return ""
	}

	buf := &bytes.Buffer{}
	for name, value := range pt {
		tag := fmt.Sprintf(`%s:%q`, name, value)
		buf.WriteString(" ")
		buf.WriteString(tag)
	}

	tagStr := buf.String()
	return tagStr[1:]
}
