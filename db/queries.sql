	-- name: GetAllTableNames :many
SELECT table_name
FROM information_schema.tables
WHERE table_schema = 'public';

-- name: ListUsers :many
SELECT name, age
FROM users
ORDER BY name;

-- name: CreateDefaultUsers :exec
INSERT INTO users (name, age)
VALUES 
    ('Alice', 30),
    ('Bob', 25),
    ('Charlie', 35),
    ('Diana', 28),
    ('Eve', 22);