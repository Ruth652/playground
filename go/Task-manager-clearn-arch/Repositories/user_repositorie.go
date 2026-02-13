package repositories

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	domain "task-manager-clean-arch/Domain"
)

type UserRepository struct {
	db  *mongo.Database
	col string
}

func NewUserRepository(db *mongo.Database, collection string) domain.UserRepository {
	return &UserRepository{
		db:  db,
		col: collection,
	}
}

func (ur *UserRepository) Create(ctx context.Context, u *domain.User) error {
	collection := ur.db.Collection(ur.col)

	_, err := collection.InsertOne(ctx, u)
	if mongo.IsDuplicateKeyError(err) {
		return domain.ErrConflict
	}
	return err
}

func (ur *UserRepository) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	collection := ur.db.Collection(ur.col)
	var user domain.User

	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return domain.User{}, domain.ErrNotFound
	}
	return user, err
}

func (ur *UserRepository) GetByID(ctx context.Context, id primitive.ObjectID) (domain.User, error) {
	collection := ur.db.Collection(ur.col)
	var user domain.User

	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return domain.User{}, domain.ErrNotFound
	}
	return user, err
}

func (ur *UserRepository) Count(ctx context.Context) (int64, error) {
	collection := ur.db.Collection(ur.col)
	return collection.CountDocuments(ctx, bson.D{})
}

func (ur *UserRepository) PromoteToAdmin(ctx context.Context, id primitive.ObjectID) error {
	collection := ur.db.Collection(ur.col)

	result, err := collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"role": domain.RoleAdmin}})
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return domain.ErrNotFound
	}
	return nil
}
