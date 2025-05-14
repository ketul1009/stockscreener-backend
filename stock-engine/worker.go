package engine

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func StartWorker(redisClient *redis.Client, stocks []Stock) {
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

		matched := FilterStocks(stocks, job.Rules)

		resultJSON, _ := json.Marshal(matched)
		err = redisClient.Set(ctx, "screener_result:"+job.JobID, resultJSON, 10*time.Minute).Err()
		if err != nil {
			fmt.Println("Failed to store result:", err)
		} else {
			fmt.Printf("✅ Job %s processed (%d matches)\n", job.JobID, len(matched))
		}
	}
}
