package redis

import (
	"context"
	"encoding/json"

	//er "errors"
	"fmt"

	"time"

	"github.com/Creative-genius001/Stacklo/services/user/model"
	errors "github.com/Creative-genius001/Stacklo/services/user/utils/error"
	"github.com/Creative-genius001/Stacklo/services/user/utils/logger"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type Redis interface {
	SaveOTPToRedis(email string, otp string) error
	GetOTPFromRedis(email string) (*model.OTPJSON, error)
	IncrementRetries(email string) error
}

type redisClient struct {
	redis *redis.Client
}

func NewRedisClient(u string) Redis {
	opt, err := redis.ParseURL(u)
	if err != nil {
		logger.Logger.Panic("Error parsing url", zap.Error(err))
	}

	return &redisClient{redis.NewClient(opt)}
}

func (r *redisClient) SaveOTPToRedis(email string, otp string) error {
	key := fmt.Sprintf("otp:%s", email)
	data := model.OTPJSON{
		OTP:       otp,
		Retry:     0,
		ExpiresAt: time.Now().Add(3 * time.Minute),
	}
	jsonData, _ := json.Marshal(data)

	cmd := r.redis.Do(context.Background(),
		"JSON.SET", key, "$", string(jsonData),
	)
	if cmd.Err() != nil {
		logger.Logger.Error("Error saving OTP to Redis", zap.Error(cmd.Err()))
		return errors.Wrap(errors.TypeInternal, "Unable to save OTP to redis", cmd.Err())
	}

	r.redis.Expire(context.Background(), key, 5*time.Minute)
	return nil
}

func (r *redisClient) GetOTPFromRedis(email string) (*model.OTPJSON, error) {
	key := fmt.Sprintf("otp:%s", email)
	val, err := r.redis.Do(context.Background(),
		"JSON.GET", key, "$",
	).Result()

	if err.Error() == "redis: nil" {
		return nil, errors.Wrap(errors.TypeInvalidInput, "code is expired", err)
	} else if err != nil {
		logger.Logger.Error("Error getting OTP from Redis", zap.Error(err))
		return nil, errors.Wrap(errors.TypeInvalidInput, "verification failed", err)
	}

	jsonD, _ := json.Marshal(val)

	var data []model.OTPJSON
	err = json.Unmarshal([]byte(jsonD), &data)
	if err != nil {
		logger.Logger.Error("invalid otp fomrat", zap.Error(err))
		return nil, err
	}

	return &data[0], nil
}

func (r *redisClient) IncrementRetries(email string) error {
	key := fmt.Sprintf("otp:%s", email)
	return r.redis.Do(context.Background(),
		"JSON.NUMINCRBY", key, "$.retry", 1,
	).Err()

}
