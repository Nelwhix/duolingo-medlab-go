package request

type QuestionType string

const (
	MultipleChoice QuestionType = "multiple_choice"
	German         QuestionType = "german"
)

type CreateQuestion struct {
	Topic         string       `validate:"required"`
	Type          QuestionType `validate:"required"`
	DepartmentID  string       `validate:"required"`
	Question      string       `validate:"required"`
	Options       []string
	CorrectOption int
	Answer        string
}
