package controllers

import (
	"net/http"

	"github.com/haitao-sun03/go/config"

	"github.com/gin-gonic/gin"
)

type TaskController struct{}

type SearchTask struct {
	UserId   int    `json:"userId" binding:"required"`
	TaskItem string `json:"taskItem"`
}

func (o TaskController) GetList(ctx *gin.Context) {
	var tasks []Task
	var total int64

	search := SearchTask{}
	if err := ctx.Bind(&search); err != nil {
		// 处理绑定错误，通常返回400 Bad Request
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := config.DB.Model(&Task{})

	if search.TaskItem != "" {
		db = db.Where("task_item LIKE ?", "%"+search.TaskItem+"%")
	}

	result := db.Where("user_id = ?", search.UserId).Find(&tasks)

	if result.Error != nil {
		panic(result.Error)
	}

	Success(ctx, http.StatusOK, "success", tasks, total)
}
