package dto

import (
	"time"

	"gitlab.strale.io/go-travel/internal/database"
)

type CommentPosterDto struct {
	ID   int64   `json:"id"`
	Name *string `json:"name,omitempty"`
}

type CommentCityDto struct {
	ID   int64   `json:"id"`
	Name *string `json:"name,omitempty"`
}

type CommentDto struct {
	ID       int64            `json:"id"`
	City     CommentCityDto   `json:"city"`
	Poster   CommentPosterDto `json:"poster"`
	Text     string           `json:"text"`
	Created  string           `json:"created,omitempty"`
	Modified string           `json:"modified,omitempty"`
}

type SaveCommentDto struct {
	CityID   int64  `json:"cityId"`
	PosterID int64  `json:"posterId"`
	Text     string `json:"text"`
}

type UpdateCommentDto struct {
	Text string `json:"text"`
}

func CommentToDto(comment database.Comment) CommentDto {
	poster := CommentPosterDto{
		ID: comment.PosterID,
	}
	if comment.Poster.ID != 0 {
		poster.Name = &comment.Poster.Username
	}
	city := CommentCityDto{
		ID: comment.CityID,
	}
	if comment.City.ID != 0 {
		city.Name = &comment.City.Name
	}
	dto := CommentDto{
		ID:     comment.ID,
		City:   city,
		Poster: poster,
		Text:   comment.Text,
	}
	if comment.Created != nil {
		dto.Created = comment.Created.Format(time.RFC3339)
	}
	if comment.Modified != nil {
		dto.Modified = comment.Modified.Format(time.RFC3339)
	}
	return dto
}

func CommentsToDtos(comments []database.Comment) []CommentDto {
	dtos := make([]CommentDto, len(comments))
	for i, comment := range comments {
		dtos[i] = CommentToDto(comment)
	}
	return dtos
}
