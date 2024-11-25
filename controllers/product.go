package controllers

import (
	"net/http"
	"time"

	"github.com/haitao-sun03/go/config"

	"github.com/gin-gonic/gin"
	"github.com/haitao-sun03/golang-distributelock/distribute"
	log "github.com/sirupsen/logrus"
)

type ProductController struct{}

func (o ProductController) SetString(ctx *gin.Context) {
	key := ctx.Param("key")
	value := ctx.Param("value")
	err := config.RedisClient.Set(ctx, key, value, 0).Err()
	if err != nil {
		// 这里应该处理错误，而不是直接panic
		Fail(ctx, http.StatusInternalServerError, err)

	}
	Success(ctx, http.StatusOK, "success", 0, 0)

}

func (o ProductController) GetString(ctx *gin.Context) {
	key := ctx.Param("key")
	// 获取键值
	val, err := config.RedisClient.Get(ctx, key).Result()
	if err != nil {
		Fail(ctx, http.StatusInternalServerError, err)
	}
	Success(ctx, http.StatusOK, "success", val, 0)
}

func (o ProductController) TestDistributeLock(ctx *gin.Context) {

	lock := distribute.NewDistributedLock(config.RedisClient, "my_lock", "lock_value", 5*time.Second)

	if acquired, err := lock.TryLock(ctx); err == nil {
		if acquired {
			log.Info("Lock acquired")
			defer func() {
				if err := lock.Unlock(ctx); err != nil {
					log.Infof("Failed to unlock: %v\n", err)

				} else {
					log.Info("Lock released")
				}
			}()
			// 在这里执行需要锁保护的操作
			time.Sleep(2 * time.Second) // 模拟一些工作
		} else {
			log.Info("Failed to acquire lock")

		}
	} else {
		log.Infof("Error acquiring lock: %v\n", err)

	}
}
