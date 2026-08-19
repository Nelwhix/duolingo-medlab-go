package models

import (
	"context"
	"time"

	"github.com/Nelwhix/duolingo-medlab-go/pkg/request"
	"github.com/oklog/ulid/v2"
)

type Topic struct {
	ID         string
	Name       string
	Department Department
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (m *Model) InsertTopic(ctx context.Context, req request.CreateTopic) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	topicID := ulid.Make().String()
	query := "insert into topics (id, name, department_id) values ($1, $2, $3)"
	_, err := m.Conn.Exec(ctx, query, topicID, req.Name, req.DepartmentID)
	if err != nil {
		return err
	}

	return nil
}

func (m *Model) GetTopics(ctx context.Context) ([]Topic, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	query := "select t.id, t.name, t.created_at, t.updated_at, d.id, d.title from topics t join departments d on d.id = t.department_id order by t.created_at desc"
	rows, err := m.Conn.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var topics []Topic
	for rows.Next() {
		var t Topic
		err := rows.Scan(&t.ID, &t.Name, &t.CreatedAt, &t.UpdatedAt, &t.Department.ID, &t.Department.Title)
		if err != nil {
			return nil, err
		}

		topics = append(topics, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return topics, nil
}
