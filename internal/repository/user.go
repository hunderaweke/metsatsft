package repository

import (
	"context"
	"log"

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

func NewUserRepository(db mongoifc.Database, ctx context.Context) domain.UserRepository {

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
		log.Fatal(err)
	}
	return &userRepository{collection: collection, ctx: ctx}
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
	for res.Next(r.ctx) {
		var user domain.User
		if err := res.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *userRepository) UpdateUser(user domain.User) (domain.User, error) {
	existingUser, err := r.GetUserByID(user.ID)
	if err != nil {
		return domain.User{}, err
	}
	if user.FirstName != existingUser.FirstName {
		existingUser.FirstName = user.FirstName
	}
	if user.LastName != existingUser.LastName {
		existingUser.LastName = user.LastName
	}
	if user.TelegramUsername != existingUser.TelegramUsername {
		existingUser.TelegramUsername = user.TelegramUsername
	}
	if user.PhoneNumber != existingUser.PhoneNumber {
		existingUser.PhoneNumber = user.PhoneNumber
	}
	objID, err := primitive.ObjectIDFromHex(user.ID)
	if err != nil {
		return domain.User{}, err
	}
	_, err = r.collection.UpdateByID(r.ctx, objID, existingUser, options.Update())
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
