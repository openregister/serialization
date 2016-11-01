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
	var fields map[string]Field
	json.Unmarshal(streamToBytes(rc), &fields)
	return fields
}

func streamToBytes(stream io.Reader) []byte {
	buf := new(bytes.Buffer)
	buf.ReadFrom(stream)
	return buf.Bytes()
}

func escapeQuotes(s string) string {
	return strings.Replace(s, `"`, `\"`, -1)
}

func toJsonArrayOfStr(s string) string {
	return `["` + strings.Replace(s, `;`, `","`, -1) + `"]`
}

func toJsonArrayOfNum(s string) string {
	return `[` + strings.Replace(s, `;`, `,`, -1) + `]`
}

func toJsonStr(r interface{}) (string, error) {
	data, err := json.Marshal(r)
	return string(data), err
}
