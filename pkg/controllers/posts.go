package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"

	"elf-server/pkg/models"
)

type PostController struct{}

func (c *PostController) GetMany(ctx *gin.Context) {
	userID := GetAuthInfo(ctx).GetUserID()

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	rows, total := models.GetPostsByUserID(userID, limit, page)
	ctx.JSON(http.StatusOK, gin.H{"total": total, "rows": rows})
}

func (c *PostController) GetOne(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	ctx.JSON(http.StatusOK, models.GetPost(uint(id)))
}

func (c *PostController) Create(ctx *gin.Context) {
	userID := GetAuthInfo(ctx).GetUserID()

	var v models.Post
	if err := ctx.BindJSON(&v); err != nil {
		glog.Error(err.Error())
	}
	v.UserID = uint(userID)
	if e := v.Save(); e == nil {
		ctx.JSON(http.StatusCreated, v)
	} else {
		HandleError(ctx, e.Error())
	}
}

func (c *PostController) Update(ctx *gin.Context) {
	userID := GetAuthInfo(ctx).GetUserID()

	id, _ := strconv.Atoi(ctx.Param("id"))
	var v models.Post
	if err := ctx.BindJSON(&v); err != nil {
		glog.Error(err.Error())
	}
	v.ID = uint(id)
	v.UserID = uint(userID)
	if e := v.Update(); e == nil {
		ctx.JSON(http.StatusOK, v)
	} else {
		HandleError(ctx, e.Error())
	}
}

func (c *PostController) Delete(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	var v models.Post
	v.ID = uint(id)
	if e := v.Delete(); e == nil {
		ctx.JSON(http.StatusNoContent, nil)
	} else {
		HandleError(ctx, e.Error())
	}
}

func (c *PostController) GetContent(ctx *gin.Context) {
	userID := GetAuthInfo(ctx).GetUserID()

	id, _ := strconv.Atoi(ctx.Param("id"))
	ctx.JSON(http.StatusOK, models.GetPostContent(uint(id), int(userID)))
}

func (c *PostController) UpdateContent(ctx *gin.Context) {
	userID := GetAuthInfo(ctx).GetUserID()

	id, _ := strconv.Atoi(ctx.Param("id"))
	var v models.Post
	if err := ctx.BindJSON(&v); err != nil {
		glog.Error(err.Error())
	}
	if e := models.UpdatePostContent(uint(id), int(userID), v.Source); e == nil {
		ctx.JSON(http.StatusOK, v)
	} else {
		HandleError(ctx, e.Error())
	}
}

func NewPostController(r gin.IRouter) *PostController {
	c := &PostController{}
	r.Group("/posts").
		GET("", c.GetMany).
		GET("/:id", c.GetOne).
		POST("", c.Create).
		PUT("/:id", c.Update).
		DELETE("/:id", c.Delete).
		OPTIONS("", nil).
		OPTIONS("/:id", nil).
		GET("/:id/content", c.GetContent).
		PUT("/:id/content", c.UpdateContent).
		OPTIONS("/:id/content", nil)
	return c
}
