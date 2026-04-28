package request

type UpdateUser struct {
	LearningGoalPassExams        *bool  `json:"learning_goal_pass_exams" validate:"required"`
	LearningGoalRefreshKnowledge *bool  `json:"learning_goal_refresh_knowledge" validate:"required"`
	LearningGoalPracticeDaily    *bool  `json:"learning_goal_practice_daily" validate:"required"`
	DepartmentID                 string `json:"department_id" validate:"required"`
}
