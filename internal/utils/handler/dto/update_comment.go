package dto

//go:generate ffjson -noencoder $GOFILE
type UpdateCommentDto struct {
	Text string `json:"text"`
}
