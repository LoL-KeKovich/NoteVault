package dto

type NoteRequest struct {
	Name  string `json:"name,omitempty"`
	Text  string `json:"text,omitempty"`
	Color string `json:"color,omitempty"`
	Media string `json:"media,omitempty"`
	Order int    `json:"order,omitempty"`
}

type NoteResponse struct {
	Data  interface{} `json:"data,omitempty"`
	Error string      `json:"error,omitempty"`
}
