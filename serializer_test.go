package main

import (
	"encoding/csv"
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
	"strings"
	"testing"
)

func TestAlphabeticalSort(t *testing.T) {
	fmt.Println("starting tests")
	fields := []string{"ba", "ca", "a"}
	sortedIndexes := alphabeticalIndexes(fields)
	expected := []int{2, 0, 1}
	for i, index := range expected {
		if sortedIndexes[i] != index {
			t.Error(`sorted indexes for b,c,a should be 2,0,1`)
		}
	}
}

func TestBuildJson(t *testing.T) {
	fieldNames := []string{"c", "b", "a"}
	fieldValues := []string{"d;h", `e "f" g`, "1"}
	sortedIndexes := []int{2, 1, 0}
	fields := map[string]Field{
		"a": Field{"1", "string", "", "", "", ""},
		"b": Field{"1", "string", "", "", "", ""},
		"c": Field{"n", "string", "", "", "", ""},
	}
	json := buildContentJson(fieldNames, fieldValues, sortedIndexes, fields)
	expected := `{"a":"1","b":"e \"f\" g","c":["d","h"]}`
	if expected != json {
		t.Error("should build " + expected)
	}
}

func TestBuildJsonIgnoreWhitespace(t *testing.T) {
	fieldNames := []string{"c", "b", "a"}
	fieldValues := []string{"d", "  ", "1"}
	sortedIndexes := []int{2, 1, 0}
	fields := map[string]Field{
		"a": Field{"1", "string", "", "", "", ""},
		"b": Field{"1", "string", "", "", "", ""},
		"c": Field{"n", "string", "", "", "", ""},
	}
	json := buildContentJson(fieldNames, fieldValues, sortedIndexes, fields)
	expected := `{"a":"1","c":["d"]}`
	if expected != json {
		t.Error("should build " + expected)
	}
}

func TestHash(t *testing.T) {
	jsonStr := "{\"key\":\"value\"}"
	jsonBytes := []byte(jsonStr)
	hash := sha256Hex(jsonBytes)
	if "e43abcf3375244839c012f9633f95862d232a95b00d5bc7348b3098b9fed7f32" != hash {
		t.Error("should hash to zzz")
	}
}

func TestReadFields(t *testing.T) {
	fieldsFile, _ := os.Open("test-data/field-records.json")
	fields := readFieldTypes(fieldsFile)
	//	fmt.Println(fields["street"].Datatype)
	if fields["street"].Datatype != "string" {
		t.Error("street field should have datatype of string")
	}
}

func TestTimestamp(t *testing.T) {
	ts := timestamp()
	if !strings.HasSuffix(ts, "Z") {
		t.Error("timestamp should be end with Z")
	}
}

// just to understand how the go csv parser works
func TestReadCommas(t *testing.T) {
	in := `"aa",bb` // default delimited is comma
	r := csv.NewReader(strings.NewReader(in))
	r.LazyQuotes = true

	records, _ := r.Read()
	name := records[0]
	if name != "aa" {
		t.Error("surrounding quotes should not be read")
	}
}

func TestReadCommasInternal(t *testing.T) {
	in := `a "bb" a,bb`
	r := csv.NewReader(strings.NewReader(in))
	r.LazyQuotes = true

	records, _ := r.Read()
	name := records[0]
	if name != `a "bb" a` {
		t.Error("internal quotes shoul be read")
	}
}
func TestReadCommasLeading(t *testing.T) {
	in := `""aa"cc",bb`
	r := csv.NewReader(strings.NewReader(in))
	r.LazyQuotes = true

	records, _ := r.Read()
	name := records[0]
	if name != `"aa"cc` {
		t.Error("leading quotes should not be read if whole string quoted")
	}
}

func TestEscape(t *testing.T) {
	escaped := escapeForJson(`aaa "bbb" ccc`)
	if escaped != `aaa \"bbb\" ccc` {
		t.Error("should escape quotes")
	}
}

func TestEscapeBackSlash(t *testing.T) {
	escaped := escapeForJson(`aaa/ccc\aaa`)
	if escaped != `aaa/ccc\\aaa` {
		t.Error("should escape slash")
	}
}

func TestEscapeNoQuotes(t *testing.T) {
	escaped := escapeForJson(`aaa bbb ccc`)
	if escaped != `aaa bbb ccc` {
		t.Error("should do nothing if no quotes")
	}
}

func TestJsonStringArray(t *testing.T) {
	json := toJsonArrayOfStr(`aa;bb;cc`)
	if json != `["aa","bb","cc"]` {
		t.Error(`should render aa;bb;cc as ["aaa","bbb","ccc"]`)
	}
}

func TestReadYaml(t *testing.T) {
	yamlFile, _ := os.Open("test-data/country.yaml")
	var register Register
	yaml.Unmarshal(streamToBytes(yamlFile), &register)

	if register.Phase != "beta" {
		t.Error(`should read phase of country as beta`)
	}
	if register.Fields[0] != "country" {
		t.Error(`should read first field of country as country`)
	}
	if register.Fields[1] != "name" {
		t.Error(`should read second field of country as name`)
	}
}

func TestMarshalRegister(t *testing.T) {
	reg := Register{"", []string{"address"}, "alpha", "address", "office-for", "Postal address"}
	json, _ := toJsonStr(reg)
	expected := `{"fields":["address"],"phase":"alpha","register":"address","registry":"office-for","text":"Postal address"}`
	if expected != json {
		t.Error(`should write json without empty fields`)
	}
}

func TestCheckFieldNames(t *testing.T) {
	fieldNames := []string{"c", "b", "a"}
	fields := map[string]Field{
		"a": Field{"1", "string", "", "", "", ""},
		"b": Field{"1", "string", "", "", "", ""},
		"c": Field{"n", "string", "", "", "", ""},
	}
	if !mapContainsAllKeys(fields, fieldNames) {
		t.Error("should confirm field names a,b,c all valid")
	}
}

func TestCheckFieldNamesMissing(t *testing.T) {
	fieldNames := []string{"d", "b", "a"}
	fields := map[string]Field{
		"a": Field{"1", "string", "", "", "", ""},
		"b": Field{"1", "string", "", "", "", ""},
		"c": Field{"n", "string", "", "", "", ""},
	}
	if mapContainsAllKeys(fields, fieldNames) {
		t.Error("should find not all field names valid")
	}
}

func TestGetKey(t *testing.T) {
	fieldNames := []string{"school", "foo"}
	fieldValues := []string{"schoolId", "bar"}
	key, _ := getKey(fieldNames, fieldValues, "school")
	if key != "schoolId" {
		t.Error("should find schoolId as value of register name field")
	}
}

func TestKeyNotFound(t *testing.T) {
	fieldNames := []string{"school", "foo"}
	fieldValues := []string{"schoolId", "bar"}
	_, err := getKey(fieldNames, fieldValues, "schoolz")
	if err.Error() != "failed to find field matching register name" {
		t.Error("should report key not found")
	}
}
