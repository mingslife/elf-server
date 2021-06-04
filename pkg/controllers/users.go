package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"

	"elf-server/pkg/models"
)

type UserController struct{}

func (c *UserController) GetMany(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	rows, total := models.GetUsers(limit, page)
	ctx.JSON(http.StatusOK, gin.H{"rows": rows, "total": total})
}

func (c *UserController) GetOne(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	ctx.JSON(http.StatusOK, models.GetUser(uint(id)))
}

func (c *UserController) Create(ctx *gin.Context) {
	var v models.User
	if err := ctx.BindJSON(&v); err != nil {
		glog.Error(err.Error())
	}
	v.Role = models.UserRoleAuthor // only allowed to create author user
	if e := v.Save(); e == nil {
		ctx.JSON(http.StatusCreated, v)
	} else {
		HandleError(ctx, e.Error())
	}
}

func (c *UserController) Update(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	var v models.User
	if err := ctx.BindJSON(&v); err != nil {
		glog.Error(err.Error())
	}
	v.ID = uint(id)
	if e := v.Update(); e == nil {
		ctx.JSON(http.StatusOK, v)
	} else {
		HandleError(ctx, e.Error())
	}
}

func (c *UserController) Delete(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	v := models.GetUser(uint(id))
	// only allowed to delete author user
	if v.Role == models.UserRoleAdmin {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "cannot delete admin user",
		})
		return
	}

	if e := v.Delete(); e == nil {
		ctx.JSON(http.StatusNoContent, nil)
	} else {
		HandleError(ctx, e.Error())
	}
}

func NewUserController(r gin.IRouter) *UserController {
	c := &UserController{}
	r.Group("/users").
		GET("", c.GetMany).
		GET("/:id", c.GetOne).
		POST("", c.Create).
		PUT("/:id", c.Update).
		DELETE("/:id", c.Delete).
		OPTIONS("", nil).
		OPTIONS("/:id", nil)
	return c
}
