package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"

	"elf-server/pkg/models"
)

type CommentController struct{}

func (c *CommentController) GetMany(ctx *gin.Context) {
	postID, _ := strconv.Atoi(ctx.Query("postId"))
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	rows, total := models.GetCommentsByPostID(postID, limit, page)
	ctx.JSON(http.StatusOK, gin.H{"rows": rows, "total": total})
}

func (c *CommentController) GetOne(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	ctx.JSON(http.StatusOK, models.GetComment(uint(id)))
}

func (c *CommentController) Create(ctx *gin.Context) {
	userID := GetAuthInfo(ctx).GetUserID()
	user := models.GetUser(userID)

	ip := ctx.ClientIP()
	userAgent := ctx.GetHeader("User-Agent")

	var v models.Comment
	if err := ctx.BindJSON(&v); err != nil {
		glog.Error(err.Error())
	}
	v.Nickname = user.Nickname
	v.Email = user.Email
	v.IP = ip
	v.UserAgent = userAgent
	v.UserID = &userID
	if e := v.Save(); e == nil {
		ctx.JSON(http.StatusCreated, v)
	} else {
		HandleError(ctx, e.Error())
	}
}

func (c *CommentController) Update(ctx *gin.Context) {
	ip := ctx.ClientIP()
	userAgent := ctx.GetHeader("User-Agent")

	id, _ := strconv.Atoi(ctx.Param("id"))
	var v models.Comment
	if err := ctx.BindJSON(&v); err != nil {
		glog.Error(err.Error())
	}
	v.ID = uint(id)
	v.IP = ip
	v.UserAgent = userAgent
	v.CommentedAt = time.Now()
	if e := v.Update(); e == nil {
		ctx.JSON(http.StatusOK, v)
	} else {
		HandleError(ctx, e.Error())
	}
}

func (c *CommentController) Delete(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	var v models.Comment
	v.ID = uint(id)
	if e := v.Delete(); e == nil {
		ctx.JSON(http.StatusNoContent, nil)
	} else {
		HandleError(ctx, e.Error())
	}
}

func (c *CommentController) Block(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	var v models.Comment
	if err := ctx.BindJSON(&v); err != nil {
		glog.Error(err.Error())
	}
	if e := models.SetCommentIsBlocked(uint(id), v.IsBlocked); e == nil {
		ctx.JSON(http.StatusOK, v)
	} else {
		HandleError(ctx, e.Error())
	}
}

func NewCommentController(r gin.IRouter) *CommentController {
	c := &CommentController{}
	r.Group("/comments").
		GET("", c.GetMany).
		GET("/:id", c.GetOne).
		POST("", c.Create).
		PUT("/:id", c.Update).
		DELETE("/:id", c.Delete).
		PUT("/:id/block", c.Block).
		OPTIONS("/:id/block", nil).
		OPTIONS("", nil).
		OPTIONS("/:id", nil)
	return c
}
