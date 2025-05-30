package data_structures

type Document struct {
	ID         string   `json:"id"`
	Content    string   `json:"content"`
	Signatures []string `json:"signatures,omitempty"`
	Approved   bool     `json:"approved,omitempty"`
}
