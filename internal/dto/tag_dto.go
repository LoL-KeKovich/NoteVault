package dto

type TagRequest struct {
	Name  string `json:"name,omitempty"`
	Color string `json:"color,omitempty"`
}

type TagResponse struct {
	Data  interface{} `json:"data,omitempty"`
	Error string      `json:"error,omitempty"`
}
