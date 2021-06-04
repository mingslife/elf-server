package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"

	"elf-server/pkg/models"
)

type NavigationController struct{}

func (c *NavigationController) GetMany(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	rows, total := models.GetNavigations(limit, page)
	ctx.JSON(http.StatusOK, gin.H{"rows": rows, "total": total})
}

func (c *NavigationController) GetOne(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	ctx.JSON(http.StatusOK, models.GetNavigation(uint(id)))
}

func (c *NavigationController) Create(ctx *gin.Context) {
	var v models.Navigation
	if err := ctx.BindJSON(&v); err != nil {
		glog.Error(err.Error())
	}
	if e := v.Save(); e == nil {
		ctx.JSON(http.StatusCreated, v)
	} else {
		HandleError(ctx, e.Error())
	}
}

func (c *NavigationController) Update(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	var v models.Navigation
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

func (c *NavigationController) Delete(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	var v models.Navigation
	v.ID = uint(id)
	if e := v.Delete(); e == nil {
		ctx.JSON(http.StatusNoContent, nil)
	} else {
		HandleError(ctx, e.Error())
	}
}

func (c *NavigationController) GetAll(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, models.GetAllNavigations())
}

func NewNavigationController(r gin.IRouter) *NavigationController {
	c := &NavigationController{}
	r.Group("/navigations").
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
