package middlewares

// import (
// 	"context"
// 	"time"

// 	errors "github.com/Creative-genius001/Stacklo/services/api-gateway/internal/utils/error"
// 	"github.com/Creative-genius001/Stacklo/services/api-gateway/internal/utils/logger"
// 	"github.com/gin-gonic/gin"
// 	"github.com/go-redis/redis/v8"
// 	"go.uber.org/zap"

// 	limiter "github.com/ulule/limiter/v3"
// 	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
// 	sredis "github.com/ulule/limiter/v3/drivers/store/redis"
// )

// func IPRateLimiter(redisURL string) gin.HandlerFunc {

// 	var c *gin.Context
// 	//1000 requests per hour
// 	rate, err := limiter.NewRateFromFormatted("1000-H")
// 	if err != nil {
// 		logger.Logger.Warn("Error parsing rate format", zap.Error(err))
// 		c.AbortWithStatusJSON(errors.GetHTTPStatus(errors.TypeInternal), gin.H{"status": "error", "message": errors.TypeInternal})
// 		return nil
// 	}

// 	// Create a redis client.
// 	option, err := redis.ParseURL(redisURL)
// 	if err != nil {
// 		logger.Logger.Warn("Error parsing Redis URL for rate limiter", zap.String("redis_url", redisURL), zap.Error(err))
// 		c.AbortWithStatusJSON(errors.GetHTTPStatus(errors.TypeInternal), gin.H{"status": "error", "message": errors.TypeInternal})
// 		return nil
// 	}
// 	client := redis.NewClient(option)

// 	// Ping Redis to ensure connection is established
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	_, err = client.Ping(ctx).Result()
// 	if err != nil {
// 		logger.Logger.Fatal("Error connecting to Redis for rate limiter", zap.String("redis_url", redisURL), zap.Error(err))
// 	}
// 	logger.Logger.Info("Successfully connected to Redis for rate limiting", zap.String("redis_url", redisURL))

// 	// Create a store with the redis client.
// 	store, err := sredis.NewStoreWithOptions(client, limiter.StoreOptions{
// 		Prefix:   "rate_limiter",
// 		MaxRetry: 3,
// 	})
// 	if err != nil {
// 		logger.Logger.Warn("Error creating redis store for rate limiter", zap.Error(err))
// 		c.AbortWithStatusJSON(errors.GetHTTPStatus(errors.TypeInternal), gin.H{"status": "error", "message": errors.TypeInternal})
// 		return nil
// 	}

// 	// Create a new middleware with the limiter instance.
// 	rateLimiterMiddleware := mgin.NewMiddleware(limiter.New(store, rate),
// 		mgin.WithKeyGetter(mgin.),
// 		mgin.WithErrorHandler(func(c *gin.Context, err error) {
// 			logger.Logger.Warn("Rate limit IP exceeded", zap.Error())
// 			c.JSON(errors.GetHTTPStatus(errors.TypeTooManyRequests), gin.H{"status": "error", "message": errors.TypeTooManyRequests})
// 			c.Abort()
// 			logger.Logger.Warn("Rate limit exceeded for IP", zap.String("ip", c.ClientIP()), zap.Error(err))
// 		}),
// 	)

// 	return rateLimiterMiddleware
// }
