-- name: CreateScreener :one
INSERT INTO screeners (user_id, name, rules)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetScreener :one
SELECT * FROM screeners
WHERE id = $1;

-- name: GetScreeners :many
SELECT * FROM screeners
WHERE user_id = $1;

-- name: UpdateScreener :one
UPDATE screeners
SET name = $2, rules = $3, stock_universe = $4
WHERE id = $1
RETURNING *;

-- name: DeleteScreener :exec
DELETE FROM screeners
WHERE id = $1;

-- name: GetJobTrackerByUserID :one
SELECT * FROM job_tracker
WHERE user_id = $1;

-- name: GetJobTrackerByJobID :one
SELECT * FROM job_tracker
WHERE job_id = $1;

-- name: CreateJobTracker :one
INSERT INTO job_tracker (job_id, user_id, job_status, job_created_at, job_updated_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: UpdateJobTrackerForNewJob :one
UPDATE job_tracker
SET job_id = $2, job_status = $3, job_created_at = $4, job_updated_at = $5
WHERE user_id = $1
RETURNING *;

-- name: UpdateJobTrackerForExistingJob :one
UPDATE job_tracker
SET job_status = $2, job_updated_at = $3
WHERE job_id = $1
RETURNING *;
