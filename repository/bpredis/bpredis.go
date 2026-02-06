// Copyright (c) 2022 Institute of Software, Chinese Academy of Sciences (ISCAS)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package bpredis

import (
	"aofs/internal/env"
	"aofs/internal/log4bp"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

var g *bpRedis
var mu sync.Mutex

func init() {
	if isTestProcess() {
		return
	}
	if err := Init(); err != nil {
		log4bp.New("", gin.Mode()).
			LogW().
			Err(err).
			Str("addr", env.REDIS_URL).
			Int("db", env.REDIS_DB).
			Msg("redis init failed during package init, will retry lazily")
	}
}

func isTestProcess() bool {
	for _, arg := range os.Args {
		if strings.HasPrefix(arg, "-test.") || strings.HasPrefix(arg, "-test=") {
			return true
		}
	}
	if len(os.Args) > 0 && strings.HasSuffix(os.Args[0], ".test") {
		return true
	}
	return false
}

type BpRediser interface {
	Set(key string, value interface{}, expiration time.Duration) error
	Get(key string) (int, error)
	GetInt64(key string) (int64, error)
	Incr(key string) error
	PushNotificationMsg(msg map[string]interface{})
	GetValue(key string) ([]byte, error)
	Close()
}

func GetRedis() BpRediser {
	if g == nil {
		if err := Init(); err != nil {
			return nil
		}
	}
	return g
}

type bpRedis struct {
	Client  *redis.Client
	Logger  *log4bp.BpLogger
	ChanMsg chan map[string]interface{}
}

func (br *bpRedis) Close() {
	br.Client.Close()
}

func Init() error {
	mu.Lock()
	defer mu.Unlock()
	if g != nil {
		return nil
	}

	var err error
	log := log4bp.New("", gin.Mode())
	client := redis.NewClient(&redis.Options{
		Addr:     env.REDIS_URL,  // redis地址
		Password: env.REDIS_PASS, // redis密码，没有则留空
		DB:       env.REDIS_DB,   // 默认数据库，默认是0
	})
	for i := 0; i < 10; i++ {
		if _, err = client.Ping().Result(); err != nil {
			log.LogW().Err(err).Int("retry", i+1).Msg("failed to connect to redis")

			time.Sleep(time.Second)
		} else {
			log.LogI().Int("db", env.REDIS_DB).Str("addr", env.REDIS_URL).Msg("connected to redis")
			break
		}
	}

	if err != nil {
		log.LogE().Err(err).Int("db", env.REDIS_DB).Str("addr", env.REDIS_URL).Msg("redis init failed after retries")
		return fmt.Errorf("failed to connect redis: %w", err)
	}

	g = &bpRedis{Client: client,
		Logger:  log4bp.New("", gin.Mode()),
		ChanMsg: make(chan map[string]interface{}, 1024)}

	return nil
}
