package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"

	"elf-server/pkg/models"
)

type CategoryController struct{}

func (c *CategoryController) GetMany(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	rows, total := models.GetCategories(limit, page)
	ctx.JSON(http.StatusOK, gin.H{"rows": rows, "total": total})
}

func (c *CategoryController) GetOne(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	ctx.JSON(http.StatusOK, models.GetCategory(uint(id)))
}

func (c *CategoryController) Create(ctx *gin.Context) {
	var v models.Category
	if err := ctx.BindJSON(&v); err != nil {
		glog.Error(err.Error())
	}
	if e := v.Save(); e == nil {
		ctx.JSON(http.StatusCreated, v)
	} else {
		HandleError(ctx, e.Error())
	}
}

func (c *CategoryController) Update(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	var v models.Category
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

func (c *CategoryController) Delete(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	var v models.Category
	v.ID = uint(id)
	if e := v.Delete(); e == nil {
		ctx.JSON(http.StatusNoContent, nil)
	} else {
		HandleError(ctx, e.Error())
	}
}

func (c *CategoryController) GetAll(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, models.GetAllCategories())
}

func NewCategoryController(r gin.IRouter) *CategoryController {
	c := &CategoryController{}
	r.Group("/categories").
		GET("", c.GetMany).
		GET("/:id", func(ctx *gin.Context) {
			switch ctx.Param("id") {
			case "all":
				c.GetAll(ctx)
			default:
				c.GetOne(ctx)
			}
		}).
		POST("", c.Create).
		PUT("/:id", c.Update).
		DELETE("/:id", c.Delete).
		OPTIONS("", nil).
		OPTIONS("/:id", nil)
	return c
}
