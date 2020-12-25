package event

// PostDocument is used to post a document.
// This creates a new document or updates an existing one.
type PostDocument struct {
	Type     string      `json:"type"`
	Name     string      `json:"name"`
	Document interface{} `json:"doc"`
}

// GetDocument is ues to get a specific document
type GetDocument struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

// ListenOnChange is used to listen to changes.
type ListenOnChange struct {
	Handler   string                   `json:"handler"`
	Documents []ListenOnChangeDocument `json:"docs"`
}

// ListenOnChangeDocument is used to listen to a changes of specific
// document properties.
type ListenOnChangeDocument struct {
	Type       string   `json:"type"`
	Name       string   `json:"name"`
	Version    uint32   `json:"version"`
	Properties []string `json:"props"`
}

// Change contains documents which where affected by a change.
type Change struct {
	Handler   string     `json:"handler"`
	Documents []Document `json:"docs"`
}

// Document contains a document content an meta data.
type Document struct {
	Type     string      `json:"type"`
	Name     string      `json:"name"`
	Version  uint32      `json:"version"`
	Document interface{} `json:"doc"`
}