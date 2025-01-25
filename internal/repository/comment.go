package repository

import (
	"context"

	"github.com/hunderaweke/metsasft/internal/domain"
	"github.com/sv-tools/mongoifc"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type commentRepository struct {
	c   mongoifc.Collection
	ctx context.Context
}

func NewCommentRepository(db mongoifc.Database, ctx context.Context) domain.CommentRepository {
	collection := db.Collection(domain.CommentCollection)
	return &commentRepository{c: collection, ctx: ctx}
}

func (r *commentRepository) CreateComment(comment domain.Comment) (domain.Comment, error) {
	res, err := r.c.InsertOne(r.ctx, comment)
	if err != nil {
		return domain.Comment{}, err
	}
	comment.ID = res.InsertedID.(string)
	return comment, nil
}
func (r *commentRepository) GetComments(filter domain.CommentFilter) ([]domain.Comment, error) {
	filterMap := bson.M{}
	if filter.BlogID != "" {
		filterMap["blog_id"] = filter.BlogID
	}
	if filter.WriterID != "" {
		filterMap["writer_id"] = filter.WriterID
	}
	cursor, err := r.c.Find(r.ctx, filterMap)
	if err != nil {
		return nil, err
	}
	var comments []domain.Comment
	err = cursor.All(r.ctx, &comments)
	if err != nil {
		return nil, err
	}
	return comments, nil

}
func (r *commentRepository) UpdateComment(comment domain.Comment) (domain.Comment, error) {
	existingComment, err := r.GetCommentByID(comment.ID)
	updateMap := bson.M{"$set": bson.M{}}
	if err != nil {
		return domain.Comment{}, err
	}
	if comment.Body != "" {
		updateMap["$set"].(bson.M)["body"] = comment.Body
		existingComment.Body = comment.Body
	}
	if comment.WriterID != "" {
		updateMap["$set"].(bson.M)["writer_id"] = comment.WriterID
		existingComment.WriterID = comment.WriterID
	}
	_, err = r.c.UpdateOne(r.ctx, bson.M{"_id": comment.ID}, updateMap)
	if err != nil {
		return domain.Comment{}, err
	}
	return existingComment, nil
}
func (r *commentRepository) DeleteComment(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.c.DeleteOne(r.ctx, bson.M{"_id": objID})
	if err != nil {
		return err
	}
	return nil
}
func (r *commentRepository) GetCommentByID(id string) (domain.Comment, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.Comment{}, err
	}
	res := r.c.FindOne(r.ctx, bson.M{"_id": objID})
	if res.Err() != nil {
		return domain.Comment{}, res.Err()
	}
	var comment domain.Comment
	err = res.Decode(&comment)
	if err != nil {
		return domain.Comment{}, err
	}
	return comment, nil
}
