package repository

import (
	"context"

	"github.com/hunderaweke/metsasft/internal/domain"
	"github.com/sv-tools/mongoifc"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type userRepository struct {
	collection mongoifc.Collection
	ctx        context.Context
}

func NewUserRepository(db mongoifc.Database, ctx context.Context) (domain.UserRepository, bool, error) {
	collection := db.Collection(domain.UserCollection)
	_, err := collection.Indexes().CreateMany(
		ctx,
		[]mongo.IndexModel{
			{
				Keys:    bson.D{{"telegram_username", 1}},
				Options: options.Index().SetUnique(true),
			},
			{
				Keys:    bson.D{{"phone_number", 1}},
				Options: options.Index().SetUnique(true),
			},
			{
				Keys:    bson.D{{"email", 1}},
				Options: options.Index().SetUnique(true),
			},
		},
	)

	if err != nil {
		return nil, false, err
	}
	cnt, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, false, err
	}
	return &userRepository{collection: collection, ctx: ctx}, cnt == 0, nil
}

func (r *userRepository) CreateUser(user domain.User) (domain.User, error) {
	res, err := r.collection.InsertOne(r.ctx, user)
	if mongo.IsDuplicateKeyError(err) {
		return domain.User{}, &domain.ErrDuplicateUser{Message: err.Error()}
	}

	objId := res.InsertedID.(primitive.ObjectID)
	user.ID = objId.Hex()
	return user, nil
}

func (r *userRepository) GetUsers(filter domain.UserFilter) ([]domain.User, error) {
	filterMap := bson.M{}
	if filter.Email != "" {
		filterMap["email"] = filter.Email
	}
	if filter.PhoneNumber != "" {
		filterMap["phone_number"] = filter.PhoneNumber
	}
	if filter.TelegramUsername != "" {
		filterMap["telegram_username"] = filter.TelegramUsername
	}
	res, err := r.collection.Find(r.ctx, filterMap, options.Find())
	if err != nil {
		return nil, err
	}
	var users []domain.User
	err = res.All(r.ctx, &users)
	return users, err
}

func (r *userRepository) UpdateUser(user domain.User) (domain.User, error) {
	existingUser, err := r.GetUserByID(user.ID)
	update := bson.M{"$set": bson.M{}}
	if err != nil {
		return domain.User{}, err
	}
	if user.FirstName != existingUser.FirstName {
		update["$set"].(bson.M)["first_name"] = user.FirstName
	}
	if user.LastName != existingUser.LastName {
		update["$set"].(bson.M)["last_name"] = user.LastName
	}
	if user.TelegramUsername != existingUser.TelegramUsername {
		update["$set"].(bson.M)["telegram_username"] = user.TelegramUsername
	}
	if user.PhoneNumber != existingUser.PhoneNumber {
		update["$set"].(bson.M)["phone_number"] = user.PhoneNumber
	}
	if user.IsActive != existingUser.IsActive {
		update["$set"].(bson.M)["is_active"] = user.IsActive
	}
	if user.IsAdmin != existingUser.IsAdmin {
		update["$set"].(bson.M)["is_admin"] = user.IsAdmin
	}
	if user.Password != existingUser.Password {
		update["$set"].(bson.M)["password"] = user.Password
	}
	objID, err := primitive.ObjectIDFromHex(user.ID)
	if err != nil {
		return domain.User{}, err
	}
	filter := bson.M{"_id": objID}
	_, err = r.collection.UpdateOne(r.ctx, filter, update, nil)
	if err != nil {
		return domain.User{}, err
	}
	return existingUser, nil
}
func (r *userRepository) DeleteUser(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.collection.DeleteOne(r.ctx, bson.M{"_id": objID}, options.Delete())
	if err != nil {
		return err
	}
	return nil
}
func (r *userRepository) GetUserByID(id string) (domain.User, error) {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.User{}, err
	}
	res := r.collection.FindOne(r.ctx, bson.M{"_id": objId})
	if res.Err() != nil {
		return domain.User{}, res.Err()
	}
	var user domain.User
	if err := res.Decode(&user); err != nil {
		return domain.User{}, err
	}
	return user, nil
}
