package resource

type UserResource struct {
	ID                           string             `json:"id"`
	Email                        string             `json:"email"`
	Username                     string             `json:"username"`
	Token                        string             `json:"token,omitempty"`
	LearningGoalPassExams        bool               `json:"learning_goal_pass_exams"`
	LearningGoalRefreshKnowledge bool               `json:"learning_goal_refresh_knowledge"`
	LearningGoalPracticeDaily    bool               `json:"learning_goal_practice_daily"`
	Department                   DepartmentResource `json:"department,omitempty"`
}

type DepartmentResource struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}
