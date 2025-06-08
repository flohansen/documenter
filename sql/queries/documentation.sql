-- name: UpsertDocumentation :exec
INSERT INTO documentations (name, content)
VALUES ($1, $2)
ON CONFLICT (name)
    DO UPDATE SET content = excluded.content;
