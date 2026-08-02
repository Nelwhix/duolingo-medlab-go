CREATE TABLE IF NOT EXISTS questions (
     id CHAR(26) PRIMARY KEY,
     type VARCHAR(50) NOT NULL,
     topic_id CHAR(26) NOT NULL,
     question TEXT NOT NULL,
     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_questions_topic_id
    ON questions (topic_id);

CREATE TABLE IF NOT EXISTS question_options (
    id CHAR(26) PRIMARY KEY,
    question_id CHAR(26) NOT NULL,
    option_text TEXT NOT NULL,
    is_correct BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_question
        FOREIGN KEY(question_id)
            REFERENCES questions(id)
            ON DELETE CASCADE
);