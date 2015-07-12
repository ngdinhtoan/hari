package parser

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sync"
)

const (
	jsonRawData = "json.RawMessage"
)

var (
	alphaNumRegex   = regexp.MustCompile("[0-9A-Za-z]+")
	startByNumRegex = regexp.MustCompile("^[0-9]")

	structMutex = &sync.Mutex{}
	structData  = map[string]json.RawMessage{}
	structDone  = map[string]bool{}
)

func registerStruct(name string, data json.RawMessage) bool {
	structMutex.Lock()
	defer structMutex.Unlock()

	if _, found := structDone[name]; found {
		return false
	}

	if _, found := structData[name]; found {
		return false
	}

	structData[name] = data
	return true
}

func markStructDone(name string) (allDone bool) {
	structMutex.Lock()
	defer structMutex.Unlock()

	structDone[name] = true
	if _, found := structData[name]; found {
		delete(structData, name)
	}

	return len(structData) == 0
}

// Parse will parse a json string to a Struct instance
func Parse(structName string, jsonData []byte, rs chan *Struct, errs chan error, done chan bool) {
	structName = nameReplacer.Replace(structName)

	if registerStruct(structName, jsonData) == false {
		return
	}

	smap, err := parseJSONToMap(jsonData)
	if err != nil {
		errs <- fmt.Errorf("error when generating struct %s: %v", structName, err)
		if markStructDone(structName) == true {
			done <- true
		}
		return
	}

	s := NewStruct(structName)
	for name, value := range smap {
		ptype := getType(value)
		pname := ToCamelCase(name)
		pname = nameReplacer.Replace(pname)
		if startByNumRegex.MatchString(pname) {
			pname = "P" + pname
		}

		if ptype == jsonRawData {
			ptype = pname
			Parse(pname, value, rs, errs, done)
		}

		property := NewStructProperty(pname, ptype, nil)
		property.AddTag("json", name)

		s.AddProperty(property)
	}

	rs <- s

	if allDone := markStructDone(structName); allDone {
		done <- true
	}

	return
}

// getType return type of value
func getType(data json.RawMessage) string {
	var i int64
	if err := json.Unmarshal(data, &i); err == nil {
		return "int64"
	}

	var f float64
	if err := json.Unmarshal(data, &f); err == nil {
		return "float64"
	}

	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		return "string"
	}

	var b bool
	if err := json.Unmarshal(data, &b); err == nil {
		return "bool"
	}

	var a []json.RawMessage
	if err := json.Unmarshal(data, &a); err == nil {
		aType := jsonRawData
		if len(a) > 0 {
			aType = getType(a[0])
		}

		return "[]" + aType
	}

	return jsonRawData
}

// parseJSONToMap will parser a string to a map string -> interface.
func parseJSONToMap(data []byte) (map[string]json.RawMessage, error) {
	v := map[string]json.RawMessage{}
	err := json.Unmarshal(data, &v)

	return v, err
}
