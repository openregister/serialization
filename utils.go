package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"strings"
	"time"
)

// utilities

func sha256Hex(b []byte) string {
	hasher := sha256.New()
	hasher.Write(b)
	return hex.EncodeToString(hasher.Sum(nil))
}

func timestamp() string {
	return time.Now().UTC().Format(time.RFC3339)
}

func readFieldTypes(rc io.Reader) map[string]Field {
	var entries map[string]FieldEntry
	json.Unmarshal(streamToBytes(rc), &entries)

	var fields = map[string]Field{}
	for fieldName, entry := range entries {
		fields[fieldName] = entry.Items[0]
	}

	return fields
}

func streamToBytes(stream io.Reader) []byte {
	buf := new(bytes.Buffer)
	buf.ReadFrom(stream)
	return buf.Bytes()
}

func escapeForJson(s string) string {
	tmp := strings.Replace(s, `\`, `\\`, -1)
	tmp = strings.Replace(tmp, `"`, `\"`, -1)
	return tmp
}

func toJsonArrayOfStr(s string) string {
	return `["` + strings.Replace(s, `;`, `","`, -1) + `"]`
}

func toJsonArrayOfNum(s string) string {
	return `[` + strings.Replace(s, `;`, `,`, -1) + `]`
}

func toJsonStr(r interface{}) (string, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	err := enc.Encode(r)
	return strings.TrimSuffix(buf.String(), "\n"), err
}

func mapContainsAllKeys(fields map[string]Field, fieldNames []string) bool {
	for _, fieldName := range fieldNames {
		if _, ok := fields[fieldName]; !ok {
			return false
		}
	}
	return true
}

func stringArrayContains(arr []string, e string) bool {
	for _, s := range arr {
		if s == e {
			return true
		}
	}
	return false
}
