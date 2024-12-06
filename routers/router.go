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

	account := r.Group("/account")
	account.PUT("/create", controllers.AccountController{}.CreateAccount)
	account.PUT("/import", controllers.AccountController{}.ImportAccount)
	account.GET("/foo", controllers.AccountController{}.Foo)
	account.GET("/wallet", controllers.AccountController{}.Wallet)
	account.GET("/block", controllers.AccountController{}.BlockHeaderAndBody)
	account.GET("/transfer", controllers.AccountController{}.TransferEther)
	account.GET("/transferToken", controllers.AccountController{}.TransferToken)
	account.POST("/mint", controllers.AccountController{}.Mint)
	account.POST("/transferTokenWithABI", controllers.AccountController{}.TransferTokenWithABI)
	account.POST("/balanceOf", controllers.AccountController{}.BalanceOf)

	return r

}
