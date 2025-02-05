package repository

import (
	"context"

	"github.com/hunderaweke/metsasft/internal/domain"
	"github.com/sv-tools/mongoifc"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type blogRepository struct {
	collec mongoifc.Collection
	ctx    context.Context
}

func NewBlogRepository(db mongoifc.Database, ctx context.Context) domain.BlogRepository {
	collection := db.Collection(domain.BlogCollection)
	return &blogRepository{collec: collection, ctx: ctx}
}

func (r *blogRepository) CreateBlog(blog domain.Blog) (domain.Blog, error) {
	res, err := r.collec.InsertOne(r.ctx, blog)
	if mongo.IsDuplicateKeyError(err) {
		return domain.Blog{}, err
	}
	objId := res.InsertedID.(primitive.ObjectID)
	blog.ID = objId.Hex()
	return blog, nil
}
func (r *blogRepository) GetBlogs(filter domain.BlogFilter) ([]domain.Blog, error) {
	filterMap := bson.M{}
	if filter.WriterID != "" {
		filterMap["writer_id"] = filter.WriterID
	}
	if filter.Status != "" {
		filterMap["status"] = filter.Status
	}
	if !filter.ModificationDateRange.StartDate.IsZero() {
		filterMap["last_modified_date"] = bson.M{
			"$gte": filter.ModificationDateRange.StartDate,
		}
	}
	if !filter.ModificationDateRange.EndDate.IsZero() {
		if filterMap["last_modified_date"] == nil {
			filterMap["last_modified_date"] = bson.M{}
		}
		filterMap["last_modified_date"].(bson.M)["$lte"] = filter.ModificationDateRange.EndDate
	}
	cursor, err := r.collec.Find(r.ctx, filterMap)
	var blogs []domain.Blog
	if err != nil {
		return nil, err
	}
	err = cursor.All(r.ctx, &blogs)
	if err != nil {
		return nil, err
	}
	return blogs, nil
}
func (r *blogRepository) UpdateBlog(blog domain.Blog) (domain.Blog, error) {
	existingBlog, err := r.GetBlogByID(blog.ID)
	if err != nil {
		return domain.Blog{}, err
	}
	updateMap := bson.M{"$set": bson.M{}}
	if blog.Title != existingBlog.Title {
		updateMap["$set"].(bson.M)["title"] = existingBlog.Title
		existingBlog.Title = blog.Title
	}
	if blog.Body != existingBlog.Body {
		updateMap["$set"].(bson.M)["body"] = existingBlog.Body
		existingBlog.Body = blog.Body
	}
	if blog.Status != existingBlog.Status {
		updateMap["$set"].(bson.M)["status"] = existingBlog.Status
		existingBlog.Status = blog.Status
	}
	if blog.LastModifiedDate.IsZero() {
		updateMap["$set"].(bson.M)["last_modified_date"] = existingBlog.LastModifiedDate
		existingBlog.LastModifiedDate = blog.LastModifiedDate
	}
	if len(updateMap["$set"].(bson.M)) == 0 {
		return existingBlog, nil
	}
	_, err = r.collec.UpdateOne(r.ctx, bson.M{"_id": blog.ID}, updateMap)
	if err != nil {
		return domain.Blog{}, err
	}
	return existingBlog, nil
}
func (r *blogRepository) DeleteBlog(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.collec.DeleteOne(r.ctx, bson.M{"_id": objID})
	if err != nil {
		return err
	}
	return nil
}
func (r *blogRepository) GetBlogByID(id string) (domain.Blog, error) {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.Blog{}, err
	}
	res := r.collec.FindOne(r.ctx, bson.M{"_id": objId})
	if res.Err() != nil {
		return domain.Blog{}, res.Err()
	}
	var blog domain.Blog
	if err := res.Decode(&blog); err != nil {
		return domain.Blog{}, err
	}
	return blog, nil
}
