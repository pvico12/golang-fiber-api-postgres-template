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