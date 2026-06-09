package model

import "time"

type Task struct {
	ID        string    `json:"id"`
	Title     string    `json:"title" firestore:"title"`
	Done      bool      `json:"done" firestore:"done"`
	CreatedAt time.Time `json:"created_at" firestore:"created_at"`
}
