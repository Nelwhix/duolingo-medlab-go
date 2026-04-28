ALTER TABLE users
ADD COLUMN learning_goal_pass_exams BOOLEAN NOT NULL DEFAULT FALSE,
ADD COLUMN learning_goal_refresh_knowledge BOOLEAN NOT NULL DEFAULT FALSE,
ADD COLUMN learning_goal_practice_daily BOOLEAN NOT NULL DEFAULT FALSE;
