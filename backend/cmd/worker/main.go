package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"
)

// Async worker (Arnela parity): BRPOP from a Redis list. Enqueue with LPUSH from the API or ops tooling.
func main() {
	addr := strings.TrimSpace(os.Getenv("REDIS_URL"))
	if addr == "" {
		addr = "localhost:6379"
	}
	queue := strings.TrimSpace(os.Getenv("WORKER_QUEUE_KEY"))
	if queue == "" {
		queue = "gonsgarage:jobs"
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:         addr,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  0,
		WriteTimeout: 3 * time.Second,
	})
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("redis ping: %v", err)
	}
	log.Printf("worker: Redis OK, BRPOP %s (timeout 5s)", queue)

	for {
		select {
		case <-ctx.Done():
			log.Printf("worker: shutdown")
			_ = rdb.Close()
			return
		default:
		}

		res, err := rdb.BRPop(ctx, 5*time.Second, queue).Result()
		if err == redis.Nil {
			continue
		}
		if err != nil {
			if ctx.Err() != nil {
				_ = rdb.Close()
				return
			}
			log.Printf("worker: BRPop error: %v", err)
			time.Sleep(time.Second)
			continue
		}
		if len(res) >= 2 {
			log.Printf("worker: job payload=%s", res[1])
		}
	}
}
