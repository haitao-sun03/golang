package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type TestController struct {
	LargeArray [10000]int // 增加一个大数组
}

var slice = make([]TestController, 10)

func (TestController) Test(ctx *gin.Context) {

	for {
		// fmt.Println("hello")
		slice = append(slice, TestController{})
	}

	Success(ctx, http.StatusOK, "success", 0, 0)
}
