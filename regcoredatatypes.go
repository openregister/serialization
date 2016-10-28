package main

type Datatype struct {
	Datatype string `json:"datatype,omitempty"`
	Phase    string `json:"phase,omitempty"`
	Text     string `json:"text,omitempty"`
}

type Field struct {
	Cardinality string `json:"cardinality,omitempty"`
	Datatype    string `json:"datatype,omitempty"`
	Field       string `json:"field,omitempty"`
	Phase       string `json:"phase,omitempty"`
	Register    string `json:"register,omitempty"`
	Text        string `json:"text,omitempty"`
}

type Register struct {
	Copyright string   `json:"copyright,omitempty"`
	Fields    []string `json:"fields,omitempty"`
	Phase     string   `json:"phase,omitempty"`
	Register  string   `json:"register,omitempty"`
	Registry  string   `json:"registry,omitempty"`
	Text      string   `json:"text,omitempty"`
}

type Registry struct {
	Registry string `json:"registry,omitempty"`
}

// need to implement sort.Interface for FieldIndex
type FieldIndex struct {
	Field string
	Index int
}

type ByAlphabetical []FieldIndex

func (a ByAlphabetical) Len() int           { return len(a) }
func (a ByAlphabetical) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByAlphabetical) Less(i, j int) bool { return a[i].Field < a[j].Field }
