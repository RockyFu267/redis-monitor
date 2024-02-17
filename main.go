package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"gopkg.in/yaml.v3"
)

// Config 结构用于解析配置文件
type Config struct {
	RedisCluster struct {
		Addresses []string `yaml:"addresses"`
		Password  string   `yaml:"password"`
	} `yaml:"redis_cluster"`
}

func main() {
	// 读取配置文件
	config, err := loadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// 创建 Redis 集群客户端
	clusterClient := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    config.RedisCluster.Addresses,
		Password: config.RedisCluster.Password,
	})

	// 使用 context 控制超时
	ctx := context.Background()

	// 连接到 Redis 集群
	if err := clusterClient.Ping(ctx).Err(); err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}

	// 每秒更新一次名为 "fuaotest" 的键的值为当前时间
	go func() {
		for {
			now := time.Now().Format(time.RFC3339)
			err := clusterClient.Set(ctx, "fuaotest", now, 0).Err()
			if err != nil {
				log.Printf("Error updating key 'fuaotest': %v", err)
			} else {
				log.Printf("Updated key 'fuaotest' with value: %s", now)
			}
			time.Sleep(time.Second)
		}
	}()

	// 每秒读取名为 "fuaotest" 的键的值
	go func() {
		for {
			val, err := clusterClient.Get(ctx, "fuaotest").Result()
			if err != nil {
				log.Printf("Error reading key 'fuaotest': %v", err)
			} else {
				log.Printf("Value of key 'fuaotest' is: %s", val)
			}
			time.Sleep(time.Second)
		}
	}()

	// 保持程序运行
	select {}
}

// loadConfig 从文件加载配置信息
func loadConfig(filename string) (*Config, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %v", err)
	}

	return &config, nil
}
