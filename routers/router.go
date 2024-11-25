package routers

import (
	"github.com/haitao-sun03/go/controllers"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// 定义中间
func MiddleWare() gin.HandlerFunc {

	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				if errs, ok := err.(error); ok {
					log.WithError(errs).Error("global error occurred")
				} else {
					log.WithField("panic", errs).Error("global error occurred")
				}

			}
		}()
		c.Next()
	}
}

func Router() *gin.Engine {
	r := gin.Default()
	// r.Use(MiddleWare())
	user := r.Group("/user")
	user.POST("/list", controllers.UserController{}.GetList)
	user.PUT("/add", controllers.UserController{}.AddUser)
	user.POST("/update", controllers.UserController{}.UpdateUser)
	user.DELETE("/delete", controllers.UserController{}.Delete)

	order := r.Group("/task")
	order.POST("/list", controllers.TaskController{}.GetList)

	product := r.Group("/product")
	product.POST("/setStr/:key/:value", controllers.ProductController{}.SetString)
	product.GET("/getStr/:key", controllers.ProductController{}.GetString)
	product.GET("/lock", controllers.ProductController{}.TestDistributeLock)

	return r

}
