package models

import (
	"context"
	"time"

	"github.com/Nelwhix/duolingo-medlab-go/pkg/request"
	"github.com/oklog/ulid/v2"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID                           string
	Username                     string
	Email                        string
	Password                     string
	LearningGoalPassExams        bool
	LearningGoalRefreshKnowledge bool
	LearningGoalPracticeDaily    bool
	Department                   Department
}

func (m *Model) GetUserByEmail(ctx context.Context, email string) (User, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var user User
	row := m.Conn.QueryRow(ctx, "select id, username, email, password from users where email = $1", email)
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (m *Model) InsertIntoUsers(ctx context.Context, request request.SignUp) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	userID := ulid.Make().String()
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(request.Password), 12)
	if err != nil {
		return "", err
	}

	sql := "insert into users(id, username, email, password) values ($1, $2, $3, $4)"
	_, err = m.Conn.Exec(ctx, sql, userID, request.Username, request.Email, string(passwordHash))

	return userID, err
}

func (m *Model) GetUserById(ctx context.Context, userID string) (User, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var user User
	query := `
		SELECT 
			u.id,
			u.username,
			u.email,
			u.password,
			u.learning_goal_pass_exams,
			u.learning_goal_refresh_knowledge,
			u.learning_goal_practice_daily,
			d.id,
			d.title
		FROM users u
		LEFT JOIN departments d ON d.id = u.department_id
		WHERE u.id = $1
	`

	row := m.Conn.QueryRow(ctx, query, userID)
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.LearningGoalPassExams,
		&user.LearningGoalRefreshKnowledge,
		&user.LearningGoalPracticeDaily,
		&user.Department.ID,
		&user.Department.Title,
	)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (m *Model) GetUserByToken(ctx context.Context, token string) (User, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var user User
	query := `SELECT u.id, u.username, u.email, u.password, u.learning_goal_pass_exams, u.learning_goal_refresh_knowledge, u.learning_goal_practice_daily 
FROM users u INNER JOIN personal_access_tokens t ON u.id = t.user_id
WHERE t.token = $1 AND t.expires_at > $2`
	row := m.Conn.QueryRow(ctx, query, token, time.Now())
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.LearningGoalPassExams, &user.LearningGoalRefreshKnowledge, &user.LearningGoalPracticeDaily)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (m *Model) UpdateUser(ctx context.Context, userID string, request request.UpdateUser) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	query := "update users set learning_goal_pass_exams = $1, learning_goal_practice_daily = $2, learning_goal_refresh_knowledge = $3, department_id = $4 where id = $5"
	_, err := m.Conn.Exec(ctx, query, request.LearningGoalPassExams, request.LearningGoalPracticeDaily, request.LearningGoalRefreshKnowledge, request.DepartmentID, userID)

	return err
}
