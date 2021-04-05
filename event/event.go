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
	Type        string   `json:"type"`
	Name        string   `json:"name"`
	NewDocument bool     `json:"newDoc"`
	Version     uint32   `json:"version"`
	Properties  []string `json:"props"`
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

// RegisterHandler regsisters a new Handler
type RegisterHandler struct {
	Handler  string `json:"handler"`
	QueueURL string `json:"queueUrl"`
}

// HandlerNewDoc is used to define that a given
// handler should be called if a new document of a given
// type was created.
type HandlerNewDoc struct {
	Handler  string `json:"handler"`
	Type     string `json:"type"`
	Existing bool   `json:"existing"`
}

// AdminCmd is used for multiple different administrative commands
type AdminCmd struct {
	Cmd             string                `json:"cmd"`
	RegisterHandler *RegisterHandler      `json:"regHandler,omitempty"`
	RegisterDocType *AdminRegisterDocType `json:"regDocType,omitempty"`
	RequeueHandler  *AdminRequeueHandler  `json:"requeueHandler,omitempty"`
	HandlerNewDoc   *HandlerNewDoc        `json:"handlerNewDoc,omitempty"`
}

// AdminRegisterDocType is used to register a new document type
type AdminRegisterDocType struct {
	Type string `json:"type"`
}

// AdminRequeueHandler is used to requeue a listeners of a handler
type AdminRequeueHandler struct {
	Handler string `json:"handler"`
}
