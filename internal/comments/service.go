package comments

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"gitlab.strale.io/go-travel/internal/comments/repository"
	"gitlab.strale.io/go-travel/internal/database"
	"gitlab.strale.io/go-travel/internal/utils/handler"
	"gorm.io/gorm"
)

type iCommentRepository interface {
	ListComments(pagination repository.Pagination) ([]database.Comment, error)
	FindByID(input repository.FindByIDInput) (*database.Comment, error)
	ListCommentsForUser(userID int64, pagination repository.Pagination) ([]database.Comment, error)
	ListCommentsForCity(cityID int64, pagination repository.Pagination) ([]database.Comment, error)
	Insert(comment database.Comment) (*database.Comment, error)
	Update(comment database.Comment) error
	Delete(id int64) error
}

type iCityRepository interface {
	FindByID(id int64, preload bool) (database.City, error)
}

type iUserRepository interface {
	FindByID(id int64, loadUserRoles bool) (database.User, error)
}

type commentService struct {
	cityRepo    iCityRepository
	commentRepo iCommentRepository
	userRepo    iUserRepository
}

func NewCommentService(
	commentRepo iCommentRepository,
	cityRepo iCityRepository,
	userRepo iUserRepository,
) *commentService {
	return &commentService{
		commentRepo: commentRepo,
		cityRepo:    cityRepo,
		userRepo:    userRepo,
	}
}

// Pagination object
type ListCommentsInput struct {
	limit  int
	offset int
}

func (cs *commentService) doListComments(
	ctx context.Context,
	searchFunc func() ([]database.Comment, error),
	errorLogFields map[string]interface{},
) ([]database.Comment, error) {
	list, err := searchFunc()
	if err != nil {
		logrus.WithContext(ctx).WithError(err).
			WithFields(errorLogFields).
			Error("could not list comments")
		return nil, fmt.Errorf("could not list comments: %w", err)
	}
	return list, nil
}

// List all comments with pagination support
func (cs *commentService) ListComments(ctx context.Context, pagination ListCommentsInput) ([]database.Comment, error) {
	return cs.doListComments(
		ctx,
		func() ([]database.Comment, error) {
			return cs.commentRepo.ListComments(repository.Pagination{
				Limit:  pagination.limit,
				Offset: pagination.offset,
			})
		},
		map[string]interface{}{
			"offset": pagination.offset,
			"limit":  pagination.limit,
		},
	)
}

func (cs *commentService) FindByID(ctx context.Context, id int64) (*database.Comment, error) {
	comment, err := cs.commentRepo.FindByID(repository.FindByIDInput{
		ID:       id,
		LoadUser: true,
		LoadCity: true,
	})
	switch {
	case err == gorm.ErrRecordNotFound:
		return nil, database.ErrNotFound
	case err != nil:
		return nil, fmt.Errorf("failed to read comment be ID: %w", err)
	default:
		return comment, nil
	}
}

func (cs *commentService) ListCommentsForCity(
	ctx context.Context,
	cityID int64,
	pagination ListCommentsInput,
) ([]database.Comment, error) {
	return cs.doListComments(
		ctx,
		func() ([]database.Comment, error) {
			return cs.commentRepo.ListCommentsForCity(
				cityID,
				repository.Pagination{
					Limit:  pagination.limit,
					Offset: pagination.offset,
				},
			)
		},
		map[string]interface{}{
			"cityId": cityID,
			"offset": pagination.offset,
			"limit":  pagination.limit,
		},
	)
}

func (cs *commentService) ListCommentsForUser(
	ctx context.Context,
	userID int64,
	pagination ListCommentsInput,
) ([]database.Comment, error) {
	return cs.doListComments(
		ctx,
		func() ([]database.Comment, error) {
			return cs.commentRepo.ListCommentsForUser(
				userID,
				repository.Pagination{
					Limit:  pagination.limit,
					Offset: pagination.offset,
				},
			)
		},
		map[string]interface{}{
			"userId": userID,
			"offset": pagination.offset,
			"limit":  pagination.limit,
		},
	)
}

// Saves new comment. If poster or city with given IDs don't exist
//
//	database.ErrNotFound
//
// is returned
func (cs *commentService) SaveComment(ctx context.Context, comment database.Comment) (*database.Comment, error) {
	_, err := cs.userRepo.FindByID(comment.PosterID, false)
	switch {
	case err == database.ErrNotFound:
		logrus.WithContext(ctx).WithError(err).
			Error("user does not exist")
		return nil, fmt.Errorf("user does not exist: %w", err)
	case err != nil:
		logrus.WithContext(ctx).WithError(err).
			Error("could not check if poster exists")
		return nil, fmt.Errorf("could not check if poster exists: %w", err)
	}
	_, err = cs.cityRepo.FindByID(comment.PosterID, false)
	switch {
	case err == database.ErrNotFound:
		logrus.WithContext(ctx).WithError(err).
			Error("city does not exist")
		return nil, fmt.Errorf("city does not exist: %w", err)
	case err != nil:
		logrus.WithContext(ctx).WithError(err).
			Error("could not check if city exists")
		return nil, fmt.Errorf("could not check if city exists: %w", err)
	}
	savedComment, err := cs.commentRepo.Insert(comment)
	if err != nil {
		logrus.WithContext(ctx).WithError(err).
			Error("could not save comment")
		return nil, fmt.Errorf("could not save comment: %w", err)
	}
	return savedComment, nil
}

func (cs *commentService) UpdateText(ctx context.Context, commentID, requestorID int64, text string) error {
	comment, err := cs.commentRepo.FindByID(repository.FindByIDInput{
		ID:       commentID,
		LoadUser: false,
		LoadCity: false,
	})
	switch {
	case err == database.ErrNotFound:
		logrus.WithContext(ctx).WithError(err).
			WithField("id", commentID).
			Error("comment not found")
		return err
	case err != nil:
		logrus.WithContext(ctx).WithError(err).
			WithField("id", commentID).
			Error("failed to check if comment exists")
		return fmt.Errorf("failed to check if comment exists: %w", err)
	}
	if comment.PosterID != requestorID {
		logrus.WithContext(ctx).
			WithFields(logrus.Fields{
				"commentId":   commentID,
				"posterId":    comment.PosterID,
				"requestorId": requestorID,
			}).
			Error("only poster can edit its comment")
		return handler.NewErrForbidden("only poster can edit its comment")
	}
	comment.Text = text
	err = cs.commentRepo.Update(*comment)
	if err != nil {
		logrus.WithContext(ctx).WithError(err).
			WithField("commentId", commentID).
			Error("failed to update comment")
		return fmt.Errorf("failed to update comment: %w", err)
	}
	return nil
}

func (cs *commentService) DeleteByID(ctx context.Context, commentID, requestorID int64, force bool) error {
	comment, err := cs.commentRepo.FindByID(repository.FindByIDInput{
		ID:       commentID,
		LoadUser: false,
		LoadCity: false,
	})
	switch {
	case err == database.ErrNotFound:
		logrus.WithContext(ctx).WithError(err).
			WithField("id", commentID).
			Error("comment not found")
		return err
	case err != nil:
		logrus.WithContext(ctx).WithError(err).
			WithField("id", commentID).
			Error("failed to check if comment exists")
		return fmt.Errorf("failed to check if comment exists: %w", err)
	}
	if !force && comment.PosterID != requestorID {
		logrus.WithContext(ctx).Error("only poster can delete comment")
		return handler.NewErrForbidden("only poster can delete comment")
	}
	if err = cs.commentRepo.Delete(commentID); err != nil {
		logrus.WithContext(ctx).WithError(err).
			WithField("commentId", commentID).
			Error("could not delete comment")
		return fmt.Errorf("could not delete comment: %w", err)
	}
	return nil
}
