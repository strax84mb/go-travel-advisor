package dto

//go:generate ffjson -noencoder $GOFILE
type SaveCommentDto struct {
	CityID   int64  `json:"cityId"`
	PosterID int64  `json:"posterId"`
	Text     string `json:"text"`
}
