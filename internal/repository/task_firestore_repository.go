package repository

import (
	"context"
	"errors"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"task-api/internal/model"
)

var ErrTaskNotFound = errors.New("task not found")

type TaskRepository struct {
	client *firestore.Client
}

func NewTaskFirestoreRepository(client *firestore.Client) *TaskRepository {
	return &TaskRepository{client: client}
}

func (r *TaskRepository) List(ctx context.Context) ([]model.Task, error) {
	iter := r.client.Collection("tasks").Documents(ctx)
	defer iter.Stop()

	tasks := make([]model.Task, 0)
	for {
		snap, err := iter.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return nil, err
		}

		var task model.Task
		if err := snap.DataTo(&task); err != nil {
			return nil, err
		}
		task.ID = snap.Ref.ID
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (r *TaskRepository) Create(ctx context.Context, title string) (model.Task, error) {
	doc := r.client.Collection("tasks").NewDoc()
	task := model.Task{
		ID:        doc.ID,
		Title:     title,
		Done:      false,
		CreatedAt: time.Now(),
	}

	if _, err := doc.Set(ctx, map[string]any{
		"title":      task.Title,
		"done":       task.Done,
		"created_at": task.CreatedAt,
	}); err != nil {
		return model.Task{}, err
	}

	return task, nil
}

func (r *TaskRepository) GetByID(ctx context.Context, id string) (model.Task, error) {
	snap, err := r.client.Collection("tasks").Doc(id).Get(ctx)
	if status.Code(err) == codes.NotFound {
		return model.Task{}, ErrTaskNotFound
	}
	if err != nil {
		return model.Task{}, err
	}

	var task model.Task
	if err := snap.DataTo(&task); err != nil {
		return model.Task{}, err
	}
	task.ID = snap.Ref.ID

	return task, nil
}

func (r *TaskRepository) Delete(ctx context.Context, id string) error {
	_, err := r.client.Collection("tasks").Doc(id).Delete(ctx)
	return err
}
