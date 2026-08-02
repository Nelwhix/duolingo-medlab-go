CREATE TABLE IF NOT EXISTS topics
(
    id            CHAR(26) PRIMARY KEY,
    name          VARCHAR(255)  NOT NULL,
    department_id CHAR(26)     NOT NULL,
    created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_topics_department_id
    ON topics (department_id);