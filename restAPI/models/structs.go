package models

// --------------------------------- STRUCT ---------------------------- //
type Article struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Content     string `json:"content"`
}

type JsonResponse struct {
	Type string    `json:"type"`
	Data []Article `json:"data"`
}
