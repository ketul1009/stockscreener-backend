package engine

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/ketul1009/stockscreener-backend/db"
	"github.com/redis/go-redis/v9"
)

type ApiConfig struct {
	DB *db.Queries
}

var ctx = context.Background()

func StartWorker(redisClient *redis.Client, stocks []Stock, cfg *ApiConfig) {
	for {
		result, err := redisClient.BRPop(ctx, 0, "screener_jobs").Result()
		if err != nil {
			fmt.Println("Queue error:", err)
			continue
		}

		var job ScreenerJob
		if err := json.Unmarshal([]byte(result[1]), &job); err != nil {
			fmt.Println("Failed to parse job:", err)
			continue
		}

		_, err = cfg.DB.UpdateJobTrackerForExistingJob(ctx, db.UpdateJobTrackerForExistingJobParams{
			JobID:        pgtype.UUID{Bytes: uuid.MustParse(job.JobID), Valid: true},
			JobStatus:    "running",
			JobUpdatedAt: pgtype.Timestamp{Time: time.Now(), Valid: true},
		})

		if err != nil {
			fmt.Println("Failed to update job tracker:", err)
			continue
		}

		matched := FilterStocks(stocks, job.Rules)

		resultJSON, _ := json.Marshal(matched)
		err = redisClient.Set(ctx, "screener_result:"+job.JobID, resultJSON, 10*time.Minute).Err()
		if err != nil {
			_, err2 := cfg.DB.UpdateJobTrackerForExistingJob(ctx, db.UpdateJobTrackerForExistingJobParams{
				JobID:        pgtype.UUID{Bytes: uuid.MustParse(job.JobID), Valid: true},
				JobStatus:    "failed",
				JobUpdatedAt: pgtype.Timestamp{Time: time.Now(), Valid: true},
			})

			if err2 != nil {
				fmt.Println("Failed to update job tracker:", err2)
			}

			fmt.Println("Failed to store result:", err)
		} else {
			_, err3 := cfg.DB.UpdateJobTrackerForExistingJob(ctx, db.UpdateJobTrackerForExistingJobParams{
				JobID:        pgtype.UUID{Bytes: uuid.MustParse(job.JobID), Valid: true},
				JobStatus:    "completed",
				JobUpdatedAt: pgtype.Timestamp{Time: time.Now(), Valid: true},
			})

			if err3 != nil {
				fmt.Println("Failed to update job tracker:", err3)
			}

			fmt.Printf("Job %s processed (%d matches)\n", job.JobID, len(matched))
		}
	}
}
