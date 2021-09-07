package controllers

import (
	"bytes"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"

	"elf-server/pkg/models"
	"elf-server/pkg/utils"
)

type PortalController struct{}

func (c *PortalController) Index(ctx *gin.Context) {
	settings := ctx.GetStringMapString("_settings")
	limit, _ := strconv.Atoi(settings["app.limit"])
	total := models.CountPostsForPortal("", "")
	pages := int(math.Ceil(float64(total) / float64(limit)))

	posts := models.GetPostsForPortal("", "", limit, 1)
	c.HTML(ctx, http.StatusOK, "index.jet", gin.H{
		"Posts": posts,
		"Limit": limit,
		"Page":  1,
		"Total": total,
		"Pages": pages,
	})
}

func (c *PortalController) Post(ctx *gin.Context) {
	route := ctx.Param("route")

	post := models.GetPostForPortal(route)
	if post == nil {
		c.HTML(ctx, http.StatusNotFound, "error.jet", gin.H{})
		return
	}
	if post.IsPrivate || post.Category.IsPrivate {
		post.Content = ""
	}
	go models.UpdatePostStatisticsPageView(post.UniqueID)
	c.HTML(ctx, http.StatusOK, "post.jet", gin.H{
		"Title":       post.Title,
		"Keywords":    post.Keywords,
		"Description": post.Description,
		"Post":        post,
	})
}

func (c *PortalController) Content(ctx *gin.Context) {
	uniqueID := ctx.Param("uniqueId")
	ticket := ctx.Query("ticket")

	post := models.GetPostForPortalByUniqueID(uniqueID)
	time1 := int(math.Floor(float64(time.Now().Unix()) / 60))
	time2 := time1 + 1
	ticket1 := utils.Md5(fmt.Sprintf("%s:%d", post.Password, time1))
	ticket2 := utils.Md5(fmt.Sprintf("%s:%d", post.Password, time2))

	if ticket == ticket1 || ticket == ticket2 {
		go models.UpdatePostStatisticsPageView(post.UniqueID)
		ctx.JSON(http.StatusOK, gin.H{
			"content": post.Content,
		})
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{})
	}
}

func (c *PortalController) Comments(ctx *gin.Context) {
	uniqueID := ctx.Param("uniqueId")
	post := models.GetPostForPortalByUniqueID(uniqueID)
	postID := post.ID
	if postID != 0 && post.IsCommentShown {
		comments, _ := models.GetCommentsByPostIDForPortal(postID)
		ctx.JSON(http.StatusOK, comments)
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{})
	}
}

type CommentRequest struct {
	ParentLevel *uint  `json:"parentLevel"`
	Nickname    string `json:"nickname"`
	Email       string `json:"email"`
	Content     string `json:"content"`
}

func (c *PortalController) Comment(ctx *gin.Context) {
	if !c.CaptchaVerify(ctx) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Wrong captcha!",
		})
		return
	}

	ip := ctx.ClientIP()
	userAgent := ctx.GetHeader("User-Agent")

	uniqueID := ctx.Param("uniqueId")
	post := models.GetPostForPortalByUniqueID(uniqueID)
	postID := post.ID
	if postID != 0 && post.IsCommentShown && post.IsCommentEnabled {
		var comment CommentRequest
		if err := ctx.BindJSON(&comment); err != nil {
			glog.Error(err.Error())
		}
		var parentID *uint
		if comment.ParentLevel != nil {
			parentComment := models.GetCommentByPostIDAndLevel(postID, *comment.ParentLevel)
			if parentComment != nil {
				parentID = &parentComment.ID
			}
		}
		(&models.Comment{
			PostID:    postID,
			IP:        ip,
			UserAgent: userAgent,
			ParentID:  parentID,
			Nickname:  comment.Nickname,
			Email:     comment.Email,
			Content:   comment.Content,
		}).Save()

		ctx.JSON(http.StatusOK, gin.H{})
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{})
	}
}

func (c *PortalController) User(ctx *gin.Context) {
	settings := ctx.GetStringMapString("_settings")
	username := ctx.Param("username")

	limit, _ := strconv.Atoi(settings["app.limit"])
	page, err := strconv.Atoi(ctx.Param("page"))
	if err != nil {
		page = 1
	}
	total := models.CountPostsForPortal(username, "")
	pages := int(math.Ceil(float64(total) / float64(limit)))

	user := models.GetUserByUsername(username)
	posts := models.GetPostsForPortal(username, "", limit, page)
	c.HTML(ctx, http.StatusOK, "user.jet", gin.H{
		"Title":       user.Nickname,
		"Keywords":    user.Tags,
		"Description": user.Introduction,
		"User":        user,
		"Posts":       posts,
		"Limit":       limit,
		"Page":        page,
		"Total":       total,
		"Pages":       pages,
	})
}

func (c *PortalController) Category(ctx *gin.Context) {
	settings := ctx.GetStringMapString("_settings")
	categoryRoute := ctx.Param("route")

	limit, _ := strconv.Atoi(settings["app.limit"])
	page, err := strconv.Atoi(ctx.Param("page"))
	if err != nil {
		page = 1
	}
	total := models.CountPostsForPortal("", categoryRoute)
	pages := int(math.Ceil(float64(total) / float64(limit)))

	category := models.GetCategoryByRoute(categoryRoute)
	posts := models.GetPostsForPortal("", categoryRoute, limit, page)
	c.HTML(ctx, http.StatusOK, "category.jet", gin.H{
		"Title":       category.CategoryName,
		"Keywords":    category.Keywords,
		"Description": category.Description,
		"Category":    category,
		"Posts":       posts,
		"Limit":       limit,
		"Page":        page,
		"Total":       total,
		"Pages":       pages,
	})
}

func (c *PortalController) Page(ctx *gin.Context) {
	settings := ctx.GetStringMapString("_settings")
	limit, _ := strconv.Atoi(settings["app.limit"])
	page, err := strconv.Atoi(ctx.Param("page"))
	if err != nil {
		page = 1
	}
	total := models.CountPostsForPortal("", "")
	pages := int(math.Ceil(float64(total) / float64(limit)))

	if page > pages {
		c.HTML(ctx, http.StatusBadRequest, "error.jet", gin.H{})
		return
	}

	posts := models.GetPostsForPortal("", "", limit, page)
	c.HTML(ctx, http.StatusOK, "page.jet", gin.H{
		"Posts": posts,
		"Limit": limit,
		"Page":  page,
		"Total": total,
		"Pages": pages,
	})
}

func (c *PortalController) Captcha(ctx *gin.Context) {
	var buf bytes.Buffer

	captchaID := captcha.New()
	captcha.WriteImage(&buf, captchaID, captcha.StdWidth, captcha.StdHeight)

	reader := bytes.NewReader(buf.Bytes())

	ctx.SetCookie("ELF_CAPTCHA_ID", captchaID, 600, "/", "", false, true)
	ctx.DataFromReader(http.StatusOK, int64(buf.Len()), "image/png", reader, nil)
}

func (c *PortalController) CaptchaVerify(ctx *gin.Context) bool {
	captchaID, _ := ctx.Cookie("ELF_CAPTCHA_ID")
	captchaCode := ctx.Query("captcha")
	return captcha.VerifyString(captchaID, captchaCode)
}

func (c *PortalController) Statistics(ctx *gin.Context) {
	uniqueID := ctx.Query("uniqueId")
	models.UpdatePostStatisticsPageView(uniqueID)
	ctx.JSON(http.StatusOK, gin.H{})
}

func (c *PortalController) HTML(ctx *gin.Context, code int, name string, obj gin.H) {
	settings := ctx.GetStringMapString("_settings")
	appTitle := settings["app.title"]
	obj["Settings"] = settings
	obj["Navigations"], _ = ctx.Get("_navigations")
	if title, ok := obj["Title"]; !ok {
		obj["Title"] = appTitle
	} else {
		obj["Title"] = fmt.Sprintf("%s | %s", title, appTitle)
	}
	if _, ok := obj["Keywords"]; !ok {
		obj["Keywords"] = settings["app.keywords"]
	}
	if _, ok := obj["Description"]; !ok {
		obj["Description"] = settings["app.description"]
	}
	ctx.HTML(code, name, obj)
}

func (c *PortalController) Common(ctx *gin.Context) {
	settings := models.GetAllSettingsMap()
	navigations := models.GetAllNavigationsActive()
	ctx.Set("_settings", settings)
	ctx.Set("_navigations", navigations)
	ctx.Next()
}

func NewPortalController(r gin.IRouter) *PortalController {
	c := &PortalController{}
	r.Use(c.Common)
	r.GET("/", c.Index)
	r.GET("/post/:route", c.Post)
	r.GET("/content/:uniqueId", c.Content)
	r.GET("/comment/:uniqueId", c.Comments)
	r.GET("/comment/:uniqueId/:page", c.Comments)
	r.POST("/comment/:uniqueId", c.Comment)
	r.GET("/user/:username", c.User)
	r.GET("/user/:username/:page", c.User)
	r.GET("/category/:route", c.Category)
	r.GET("/category/:route/:page", c.Category)
	r.GET("/page", c.Page)
	r.GET("/page/:page", c.Page)
	r.GET("/captcha", c.Captcha)
	r.GET("/statistics", c.Statistics)
	return c
}
