package request

type CreateTopic struct {
	Name         string `schema:"topic" validate:"required"`
	DepartmentID string `schema:"department_id" validate:"required"`
}
