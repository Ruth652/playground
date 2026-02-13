package repositories

import (
	"context"

	"task-manager-clean-arch/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type taskRepository struct {
	db  *mongo.Database
	col string
}

func NewTaskRepository(db *mongo.Database, collection string) domain.TaskRepository {
	return &taskRepository{
		db:  db,
		col: collection,
	}
}

func (tr *taskRepository) Create(ctx context.Context, t *domain.Task) error {
	collection := tr.db.Collection(tr.col)
	_, err := collection.InsertOne(ctx, t)
	return err
}

func (tr *taskRepository) FetchByUser(ctx context.Context, userID primitive.ObjectID) ([]domain.Task, error) {
	collection := tr.db.Collection(tr.col)

	var tasks []domain.Task

	cursor, err := collection.Find(ctx, bson.M{"userID": userID})
	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &tasks)
	return tasks, err
}

func (tr *taskRepository) Update(ctx context.Context, t *domain.Task) error {
	collection := tr.db.Collection(tr.col)

	result, err := collection.UpdateOne(ctx,
		bson.M{"_id": t.ID},
		bson.M{"$set": t},
	)

	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (tr *taskRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	collection := tr.db.Collection(tr.col)

	result, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return domain.ErrNotFound
	}

	return nil
}
func (tr *taskRepository) GetByID(ctx context.Context, id primitive.ObjectID) (domain.Task, error) {
	collection := tr.db.Collection(tr.col)

	var task domain.Task

	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&task)

	if err == mongo.ErrNoDocuments {
		return domain.Task{}, domain.ErrNotFound
	}

	return task, err
}
