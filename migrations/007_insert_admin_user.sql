ALTER TABLE users ADD COLUMN role VARCHAR(20) NOT NULL DEFAULT 'user';

INSERT INTO users (id, email, username, password, role) VALUES ('01KQ9G0ABM244F14KPH6EV3S32', 'nkayb762@gmail.com', 'admin', '{{.admin_password}}', 'admin');