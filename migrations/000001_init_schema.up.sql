-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- USERS
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

-- PROJECTS
CREATE TABLE projects (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL,
    description TEXT,
    owner_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT NOW()
);

-- TASKS
CREATE TABLE tasks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title TEXT NOT NULL,
    description TEXT,
    status TEXT NOT NULL CHECK (status IN ('todo', 'in_progress', 'done')),
    priority TEXT NOT NULL CHECK (priority IN ('low', 'medium', 'high')),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    assignee_id UUID REFERENCES users(id),
    due_date TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- INDEXES
CREATE INDEX idx_projects_owner_id ON projects(owner_id);
CREATE INDEX idx_tasks_project_id ON tasks(project_id);
CREATE INDEX idx_tasks_assignee_id ON tasks(assignee_id);



-- PASSWORD: password123 (bcrypt cost 12)
-- hash generated externally
INSERT INTO users (id, name, email, password)
VALUES (
    uuid_generate_v4(),
    'Test User',
    'test@example.com',
    '$2a$12$C6UzMDM.H6dfI/f/IKcEeO7j9z0p1h3z5Zk90un0nLBKBPXn1HULy'
);


-- Create project
INSERT INTO projects (id, name, description, owner_id)
SELECT uuid_generate_v4(), 'Sample Project', 'Demo project', id
FROM users WHERE email = 'test@example.com';

-- Create tasks
INSERT INTO tasks (id, title, status, priority, project_id)
SELECT uuid_generate_v4(), 'Task 1', 'todo', 'medium', id FROM projects LIMIT 1;

INSERT INTO tasks (id, title, status, priority, project_id)
SELECT uuid_generate_v4(), 'Task 2', 'in_progress', 'high', id FROM projects LIMIT 1;

INSERT INTO tasks (id, title, status, priority, project_id)
SELECT uuid_generate_v4(), 'Task 3', 'done', 'low', id FROM projects LIMIT 1;