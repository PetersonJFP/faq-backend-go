-- name: ListFAQs :many
SELECT *
FROM faqs
ORDER BY id;
-- name: CreateFAQ :one
INSERT INTO faqs (question, answer, is_premium)
VALUES ($1, $2, $3)
RETURNING *;
-- name: DeleteFAQ :exec
DELETE FROM faqs
WHERE id = $1;
-- name: UpdateFAQ :one
UPDATE faqs
SET question = $2,
  answer = $3,
  is_premium = $4
WHERE id = $1
RETURNING *;