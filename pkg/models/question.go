package models

import (
	"context"
	"time"

	"github.com/Nelwhix/duolingo-medlab-go/pkg/request"
	"github.com/oklog/ulid/v2"
)

type QuestionOption struct {
	ID         string
	QuestionID string
	OptionText string
	IsCorrect  bool
}

type Question struct {
	ID         string
	Topic      string
	Type       request.QuestionType
	Department Department
	Question   string
	Options    []QuestionOption
	CreatedAt  time.Time
}

func (m *Model) InsertQuestion(ctx context.Context, req request.CreateQuestion) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	tx, err := m.Conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	questionID := ulid.Make().String()
	query := "insert into questions (id, topic, type, department_id, question) values ($1, $2, $3, $4, $5)"
	_, err = tx.Exec(ctx, query, questionID, req.Topic, req.Type, req.DepartmentID, req.Question)
	if err != nil {
		return err
	}

	query2 := "insert into question_options (id, question_id, option_text, is_correct) values ($1, $2, $3, $4)"
	if req.Type == request.German {
		optionID := ulid.Make().String()

		_, err = tx.Exec(ctx, query2, optionID, questionID, req.Answer, true)
		if err != nil {
			return err
		}
	}

	if req.Type == request.MultipleChoice {
		for idx, option := range req.Options {
			optionID := ulid.Make().String()

			_, err = tx.Exec(ctx, query2, optionID, questionID, option, idx == req.CorrectOption)
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit(ctx)
}

func (m *Model) GetQuestions(ctx context.Context) ([]Question, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	query := `SELECT 
    q.id,
    q.topic,
    q.type,
    q.department_id,
    q.question,
    q.created_at,
    qo.id,
    qo.question_id,
    qo.option_text,
    qo.is_correct,
    d.title
FROM questions q
INNER JOIN departments d ON d.id = q.department_id
LEFT JOIN question_options qo ON qo.question_id = q.id
ORDER BY q.created_at DESC, qo.created_at;`
	rows, err := m.Conn.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	questionMap := make(map[string]*Question)
	questionOrder := make([]string, 0)

	for rows.Next() {
		var q Question
		var optionID *string
		var optionQuestionID *string
		var optionText *string
		var isCorrect *bool

		err = rows.Scan(&q.ID, &q.Topic, &q.Type, &q.Department.ID, &q.Question, &q.CreatedAt, &optionID, &optionQuestionID, &optionText, &isCorrect, &q.Department.Title)
		if err != nil {
			return nil, err
		}

		existingQuestion, ok := questionMap[q.ID]
		if !ok {
			q.Options = []QuestionOption{}
			questionMap[q.ID] = &q
			questionOrder = append(questionOrder, q.ID)
			existingQuestion = &q
		}

		existingQuestion.Options = append(existingQuestion.Options, QuestionOption{
			ID:         *optionID,
			QuestionID: *optionQuestionID,
			OptionText: *optionText,
			IsCorrect:  *isCorrect,
		})
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	questions := make([]Question, 0, len(questionOrder))
	for _, id := range questionOrder {
		questions = append(questions, *questionMap[id])
	}

	return questions, nil
}

func (m *Model) DeleteQuestion(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	_, err := m.Conn.Exec(ctx, "DELETE FROM questions WHERE id = $1", id)
	return err
}

func (m *Model) GetQuestionById(ctx context.Context, id string) (*Question, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	query := `SELECT 
        q.id,
        q.topic,
        q.type,
        q.department_id,
        q.question,
        q.created_at,
        qo.id,
        qo.question_id,
        qo.option_text,
        qo.is_correct,
        d.title
    FROM questions q
    INNER JOIN departments d ON d.id = q.department_id
    LEFT JOIN question_options qo ON qo.question_id = q.id
    WHERE q.id = $1
    ORDER BY qo.created_at;`

	rows, err := m.Conn.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var question *Question

	for rows.Next() {
		var q Question
		var optionID *string
		var optionQuestionID *string
		var optionText *string
		var isCorrect *bool

		err = rows.Scan(&q.ID, &q.Topic, &q.Type, &q.Department.ID, &q.Question, &q.CreatedAt, &optionID, &optionQuestionID, &optionText, &isCorrect, &q.Department.Title)
		if err != nil {
			return nil, err
		}

		if question == nil {
			q.Options = []QuestionOption{}
			question = &q
		}

		question.Options = append(question.Options, QuestionOption{
			ID:         *optionID,
			QuestionID: *optionQuestionID,
			OptionText: *optionText,
			IsCorrect:  *isCorrect,
		})
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if question == nil {
		return nil, nil
	}

	return question, nil
}
