package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"
	"time"
)

func buildContentJson(fieldNames []string, fieldValues []string, sortedIndexes []int, fields map[string]Field) string {
	jsonParts := []string{}
	for _, index := range sortedIndexes {
		if len(strings.TrimSpace(fieldValues[index])) > 0 {
			jsonPart := ""
			fieldDef := fields[fieldNames[index]]
			escapedValue := escapeForJson(fieldValues[index])
			if fieldDef.Cardinality == "n" {
				if fieldDef.Datatype == "integer" {
					jsonPart = fmt.Sprintf(`"%s":%s`, fieldNames[index], toJsonArrayOfNum(escapedValue))
				} else {
					jsonPart = fmt.Sprintf(`"%s":%s`, fieldNames[index], toJsonArrayOfStr(escapedValue))
				}
			} else {
				jsonPart = fmt.Sprintf(`"%s":"%s"`, fieldNames[index], escapedValue)
			}
			jsonParts = append(jsonParts, jsonPart)
		}
	}
	jsonBody := strings.Join(jsonParts, ",")
	return "{" + jsonBody + "}"
}

func alphabeticalIndexes(fields []string) []int {
	fieldIndexes := make([]FieldIndex, len(fields))
	for index, field := range fields {
		fieldIndexes[index] = FieldIndex{field, index}
	}

	sort.Sort(ByAlphabetical(fieldIndexes))

	sortedIndexes := make([]int, len(fields))
	for i, fieldIndex := range fieldIndexes {
		sortedIndexes[i] = fieldIndex.Index
	}
	return sortedIndexes
}

func getKey(fieldNames []string, fieldValues []string, registerName string) (string, error) {
	for i, fieldName := range fieldNames {
		if fieldName == registerName {
			key := fieldValues[i]
			if len(key) > 0 {
				return key, nil
			} else {
				return "", errors.New("failed to find field matching register name")
			}
		}
	}
	return "", errors.New("failed to find field matching register name")
}

func processEmptyRootHash() {
	fmt.Println("assert-root-hash\tsha-256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855")
}

func processLine(fieldValues []string, fieldNames []string, sortedIndexes []int, fieldDefns map[string]Field, registerName string) (string, string, error) {
	key, err := getKey(fieldNames, fieldValues, registerName)
	if err != nil {
		log.Fatal("Error: getting key " + err.Error())
		return "", "", err
	}
	contentJson := buildContentJson(fieldNames, fieldValues, sortedIndexes, fieldDefns)
	contentJsonHash := "sha-256:" + sha256Hex([]byte(contentJson))
	entryParts := []string{"append-entry", "user", key, timestamp(), contentJsonHash}
	entryLine := strings.Join(entryParts, "\t")
	itemParts := []string{"add-item", string(contentJson)}
	itemLine := strings.Join(itemParts, "\t")
	return itemLine, entryLine, nil
}

func processCSV(fieldsFile, tsvFile io.Reader, registerName string, includeRootHash bool) {
	var fields = map[string]Field{}
	fields, err := readFieldTypes(fieldsFile)
	if err != nil {
		log.Fatal("Error: extracting fields: " + err.Error())
		return
	}

	csvReader := csv.NewReader(tsvFile)
	csvReader.Comma = '\t'
	csvReader.LazyQuotes = true
	//read header
	fieldNames, err := csvReader.Read()
	if err != nil {
		log.Fatal("Error: reading first line of tsv file: " + err.Error())
		return
	}
	if !mapContainsAllKeys(fields, fieldNames) {
		log.Fatal("Error: fields in tsv did not match fields json: " + fmt.Sprint(fieldNames))
		return
	}
	if !stringArrayContains(fieldNames, registerName) {
		log.Fatal("Error: field headings do not include register name " + fmt.Sprint(fieldNames))
		return
	}

	if includeRootHash {
		processEmptyRootHash()
	}

	sortedIndexes := alphabeticalIndexes(fieldNames)
	for {
		fieldValues, err := csvReader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal("Error reading csv line: " + err.Error())
			return
		}
		itemLine, entryLine, err := processLine(fieldValues, fieldNames, sortedIndexes, fields, registerName)
		if err != nil {
			return
		}
		fmt.Println(itemLine)
		fmt.Println(entryLine)
	}
}

func processYamlFile(fileInfo os.FileInfo, yamlDir string, registerName string) {
	if strings.HasSuffix(fileInfo.Name(), ".yaml") {
		yamlFile, err := os.Open(yamlDir + "/" + fileInfo.Name())
		if err != nil {
			log.Fatal("Error reading yaml file: " + err.Error())
			return
		}
		defer yamlFile.Close()

		itemLine, entryLine, err := processYaml(yamlFile, registerName)
		if err != nil {
			log.Fatal(err.Error())
			return
		}
		fmt.Println(itemLine)
		fmt.Println(entryLine)
	}
}

func processYaml(yamlFile io.Reader, registerName string) (string, string, error) {
	var contentJson string
	var err error
	var key string

	switch registerName {
	case "datatype":
		var r Datatype
		yaml.Unmarshal(streamToBytes(yamlFile), &r)
		contentJson, err = toJsonStr(r)
		key = r.Datatype
	case "field":
		var r Field
		yaml.Unmarshal(streamToBytes(yamlFile), &r)
		contentJson, err = toJsonStr(r)
		key = r.Field
	case "register":
		var r Register
		yaml.Unmarshal(streamToBytes(yamlFile), &r)
		contentJson, err = toJsonStr(r)
		key = r.Register
	default:
		return "", "", errors.New("Error: register name not recognised " + registerName)
	}
	if err != nil {
		return "", "", errors.New("Error: failed to marshal to json for " + string(streamToBytes(yamlFile)))
	}

	contentJsonHash := "sha-256:" + sha256Hex([]byte(contentJson))
	entryParts := []string{"append-entry", "user", key, timestamp(), contentJsonHash}
	entryLine := strings.Join(entryParts, "\t")
	itemParts := []string{"add-item", string(contentJson)}
	itemLine := strings.Join(itemParts, "\t")
	return itemLine, entryLine, nil
}

func main() {
	if len(os.Args) < 4 {
		log.Fatal("Usage: serializer tsv|yaml <fields json file> <data file/directory> <register name> [-excludeRootHash]")
		return
	}

	log.Println(time.Now())

	includeRootHash := true
	if len(os.Args) > 5 {
		includeRootHash = os.Args[5] != "-excludeRootHash"
	}

	registerName := os.Args[4]
	fieldsFileName := os.Args[2]
	fieldsFile, fieldsErr := os.Open(fieldsFileName)
	if fieldsErr != nil {
		log.Fatal("Error: reading fields json file: " + fieldsErr.Error())
		return
	}
	defer fieldsFile.Close()

	switch os.Args[1] {

	case "tsv":
		tsvFileName := os.Args[3]
		tsvFile, err := os.Open(tsvFileName)
		if err != nil {
			log.Fatal("Error: reading tsv file: " + err.Error())
			return
		}
		defer tsvFile.Close()
		processCSV(fieldsFile, tsvFile, registerName, includeRootHash)

	case "yaml":
		yamlDir := os.Args[3]
		files, err := ioutil.ReadDir(yamlDir)
		if err != nil {
			log.Fatal("Error: reading yaml directory: " + err.Error())
			return
		}

		if includeRootHash {
			processEmptyRootHash()
		}
		for _, file := range files {
			processYamlFile(file, yamlDir, registerName)
		}

	default:
		log.Fatal("Error: file type was not 'yaml' or 'tsv'")
	}

	log.Println(time.Now())
}
