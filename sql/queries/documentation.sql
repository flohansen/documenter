-- name: UpsertDocumentation :exec
INSERT INTO documentations (name, content)
VALUES ($1, $2)
ON CONFLICT (name)
    DO UPDATE SET content = excluded.content;

-- name: GetDocumentations :many
SELECT name FROM documentations;

-- name: GetDocumentationByName :one
SELECT content FROM documentations WHERE name = $1 LIMIT 1;
