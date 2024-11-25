package controllers

import (
	"net/http"

	"github.com/haitao-sun03/go/config"

	"github.com/gin-gonic/gin"
)

type UserController struct{}

type PageUser struct {
	Pagination
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type UserIn struct {
	Id    int    `gorm:"primary_key"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Tasks []Task `json:"tasks" gorm:"foreignKey:UserId"`
}

func (userIn UserIn) TableName() string {
	return "user" // 指定表名为 'user'
}

type Task struct {
	Id       int    `json:"id" gorm:"primary_key"`
	UserId   int    `json:"userId"`
	TaskItem string `json:"taskItem"`
}

func (task Task) TableName() string {
	return "task" // 指定表名为 'task'
}

func (u UserController) AddUser(ctx *gin.Context) {
	userIn := UserIn{}
	ctx.Bind(&userIn)

	// 处理新增的 Tasks 字段
	if len(userIn.Tasks) == 0 {
		userIn.Tasks = []Task{
			{TaskItem: "default task"},
		}
	}

	db := config.DB.Model(&UserIn{})
	db.Create(&userIn)
	Success(ctx, http.StatusOK, "success", nil, 0)
}

func (u UserController) GetUserInfo(ctx *gin.Context) {
	Success(ctx, http.StatusOK, "success", "user info", 1)
}

func (u UserController) GetList(ctx *gin.Context) {
	var users []UserIn
	var total int64

	pageUser := PageUser{}

	if err := ctx.ShouldBindJSON(&pageUser); err != nil {
		// 处理绑定错误，通常返回400 Bad Request
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := config.DB.Model(&UserIn{})

	if len(pageUser.Name) != 0 {
		db = db.Where("name LIKE ?", "%"+pageUser.Name+"%")
	}

	if pageUser.Age != 0 {
		db = db.Where("age = ?", pageUser.Age)
	}
	// 计算总记录数
	db.Count(&total)

	PaginationFunc(db, pageUser.Page, pageUser.PageSize)

	// 查询分页数据
	result := db.Order("id asc").Find(&users)

	if result.Error != nil {
		panic(result.Error)
	}

	Success(ctx, http.StatusOK, "success", users, total)
}

func (u UserController) UpdateUser(ctx *gin.Context) {
	user := UserIn{}
	ctx.Bind(&user)
	newTask := user.Tasks

	dbUser := config.DB.Model(&UserIn{})
	dbTask := config.DB.Model(&Task{})

	dbUser.Where("id=?", user.Id).Update("name", user.Name).Update("age", user.Age)
	// 删除该用户之前的task
	dbTask.Where("user_id =?", user.Id).Delete(&Task{})
	dbTask.Create(&newTask)

	Success(ctx, http.StatusOK, "success", nil, 0)

}

func (u UserController) Delete(ctx *gin.Context) {
	user := UserIn{}
	ctx.Bind(&user)

	dbUser := config.DB.Model(&UserIn{})
	dbTask := config.DB.Model(&Task{})

	dbTask.Where("user_id =?", user.Id).Delete(&Task{})
	dbUser.Delete(&user)
	Success(ctx, http.StatusOK, "success", nil, 0)

}
