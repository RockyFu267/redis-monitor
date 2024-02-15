package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

func main() {
	// 创建 Redis 集群客户端
	clusterClient := redis.NewClusterClient(&redis.ClusterOptions{
		// Addrs:    []string{"redis-cluster-node1:6379", "redis-cluster-node2:6379", "redis-cluster-node3:6379"},
		Addrs:    []string{"redis-svc:6379"},
		Password: "password?",
	})

	// 使用 context 控制超时
	ctx := context.Background()

	// 连接到 Redis 集群
	if err := clusterClient.Ping(ctx).Err(); err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	} else {
		fmt.Println("connect success")
	}

	// 每秒更新一次名为 "a" 的键的值为当前时间
	go func() {
		for {
			now := time.Now().Format(time.RFC3339)
			err := clusterClient.Set(ctx, "a", now, 0).Err()
			if err != nil {
				log.Printf("Error updating key 'a': %v", err)
			} else {
				log.Printf("Updated key 'a' with value: %s", now)
			}
			time.Sleep(time.Second)
		}
	}()

	// 每秒读取名为 "a" 的键的值
	go func() {
		for {
			val, err := clusterClient.Get(ctx, "a").Result()
			if err != nil {
				log.Printf("Error reading key 'a': %v", err)
			} else {
				log.Printf("Value of key 'a' is: %s", val)
				fmt.Println(val)
			}
			time.Sleep(time.Second)
		}
	}()

	// 保持程序运行
	select {}
}
