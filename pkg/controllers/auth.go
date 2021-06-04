package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"

	"elf-server/pkg/middleware"
	"elf-server/pkg/models"
)

type AuthController struct{}

type AuthUser struct {
	Account  string `json:"account"`
	Password string `json:"password"`
}

type AuthInfo map[string]interface{}

func (m AuthInfo) GetUserID() uint {
	userID, _ := strconv.Atoi(m["sub"].(string))
	return uint(userID)
}

func (m AuthInfo) GetUserRole() uint8 {
	userRole, _ := strconv.Atoi(m["rol"].(string))
	return uint8(userRole)
}

func (m AuthInfo) GetUserData() (userID uint, userRole uint8) {
	userID = m.GetUserID()
	userRole = m.GetUserRole()
	return
}

func (c *AuthController) Login(ctx *gin.Context) {
	authUser := &AuthUser{}
	if err := ctx.BindJSON(authUser); err != nil {
		return
	}
	if user := models.GetUserByAccountAndPassword(authUser.Account, authUser.Password); user != nil {
		token, _, _ := middleware.JWT.GenerateToken(map[string]interface{}{
			"sub": fmt.Sprintf("%d", user.ID),
			"rol": fmt.Sprintf("%d", user.Role),
		})
		ctx.JSON(http.StatusOK, gin.H{
			"token": token,
		})
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "incorrect account or password"})
	}
}

func (c *AuthController) Logout(ctx *gin.Context) {}

func (c *AuthController) Refresh(ctx *gin.Context) {
	if token, _, err := middleware.JWT.RefreshToken(ctx); err == nil {
		ctx.JSON(http.StatusOK, gin.H{
			"token": token,
		})
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
	}
}

func (c *AuthController) GetProfile(ctx *gin.Context) {
	userID := GetAuthInfo(ctx).GetUserID()

	ctx.JSON(http.StatusOK, models.GetUser(userID))
}

func (c *AuthController) UpdateProfile(ctx *gin.Context) {
	userID := GetAuthInfo(ctx).GetUserID()

	var v models.User
	if err := ctx.BindJSON(&v); err != nil {
		glog.Error(err.Error())
	}
	v.ID = userID
	v.IsActive = true
	if e := v.Update(); e == nil {
		ctx.JSON(http.StatusOK, v)
	} else {
		HandleError(ctx, e.Error())
	}
}

func (c *AuthController) Register(ctx *gin.Context) {
	inviteCode := ctx.Query("inviteCode")

	settings := models.GetAllSettingsMap()
	if settings["app.register"] != "true" || inviteCode != settings["app.inviteCode"] {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Register failed"})
		return
	}

	var v models.User
	if err := ctx.ShouldBindJSON(&v); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// check for username, email and phone
	if models.ExistsUser(v.Username, v.Email, v.Phone) {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Account alreay exists"})
		return
	}

	now := time.Now()
	user := models.User{
		Username:     v.Username,
		Password:     v.Password,
		Nickname:     v.Nickname,
		Email:        v.Email,
		Phone:        v.Phone,
		Tags:         "",
		Introduction: "",
		IsActive:     true,
		ActiveAt:     &now,
		Avatar:       "/assets/avatar.svg",
		Gender:       v.Gender,
		Birthday:     v.Birthday,
		Role:         models.UserRoleAuthor,
	}
	if e := user.Save(); e == nil {
		ctx.JSON(http.StatusOK, gin.H{})
	} else {
		HandleError(ctx, e.Error())
	}
}

func (c *AuthController) Settings(ctx *gin.Context) {
	settings := models.GetAllPublicSettingsMap()
	ctx.JSON(http.StatusOK, settings)
}

func (c *AuthController) Info(ctx *gin.Context) {
	m := GetAuthInfo(ctx)
	ctx.JSON(http.StatusOK, gin.H{
		"userId":   m.GetUserID(),
		"userRole": m.GetUserRole(),
	})
}

func NewAuthController(r gin.IRouter) *AuthController {
	c := &AuthController{}
	r.Group("/auth").
		POST("/login", c.Login).
		POST("/logout", c.Logout).
		POST("/refresh", c.Refresh).
		GET("/profile", c.GetProfile).
		PUT("/profile", c.UpdateProfile).
		POST("/register", c.Register).
		GET("/settings", c.Settings).
		GET("/info", c.Info).
		OPTIONS("/:route", nil)
	return c
}

func GetAuthInfo(ctx *gin.Context) AuthInfo {
	payload, _ := middleware.JWT.GetPayload(ctx)
	return payload
}
