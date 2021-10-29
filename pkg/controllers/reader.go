package controllers

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
	"gopkg.in/gomail.v2"

	"elf-server/pkg/components"
	"elf-server/pkg/locales"
	"elf-server/pkg/models"
	"elf-server/pkg/utils"
)

const (
	ReaderKey             = "_READER"
	ReaderEmail           = "_READER_EMAIL"
	ReaderTokenCookieName = "ELF_READER_TOKEN"
)

// ReaderApiController Reader API, this module will become an independent system in some days
type ReaderApiController struct{}

func readerTokenKey(email string) string {
	return fmt.Sprintf("reader_token:%s", email)
}

func GetReaderFromCookie(ctx *gin.Context) (*models.Reader, string) {
	t, _ := ctx.Cookie(ReaderTokenCookieName)
	ts := strings.Split(t, " ")
	if len(ts) != 2 {
		return nil, ""
	}
	email, token := ts[0], ts[1]
	if email == "" || token == "" {
		return nil, ""
	}

	readerTokenKey := readerTokenKey(email)
	readerToken := components.Cache.Get(readerTokenKey)
	if token == readerToken {
		reader := models.GetReaderByEmail(email)
		return reader, email
	}

	return nil, ""
}

func (c *ReaderApiController) Middleware(ctx *gin.Context) {
	defer ctx.Next()

	reader, email := GetReaderFromCookie(ctx)

	ctx.Set(ReaderKey, reader)
	ctx.Set(ReaderEmail, email)
}

func (c *ReaderApiController) readerTokenKey(email string) string {
	return readerTokenKey(email)
}

func (c *ReaderApiController) validateCodeKey(email string) string {
	return fmt.Sprintf("reader_validate_code:%s", email)
}

type ReaderApiSendCodeReq struct {
	Email   string `json:"email"`
	Captcha string `json:"captcha"`
}

type ReaderApiSendCodeRsp struct{}

func (c *ReaderApiController) SendCode(ctx *gin.Context) {
	settings := models.GetAllSettingsMap()
	l := locales.Locales[settings["app.language"]]

	req, rsp := &ReaderApiSendCodeReq{}, &ReaderApiSendCodeRsp{}
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	captchaID, _ := ctx.Cookie("ELF_CAPTCHA_ID")
	if !captcha.VerifyString(captchaID, req.Captcha) {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": l["reader.wrong_captcha"]})
		return
	}

	validateCode := utils.RandDigists(6)

	m := gomail.NewMessage()
	m.SetHeader("From", settings["app.smtpEmail"])
	m.SetHeader("To", req.Email)
	m.SetHeader("Subject", strings.ReplaceAll(l["reader.email_subject_send_code"], "{title}", settings["app.title"]))
	m.SetBody("text/html;charset=UTF-8", strings.ReplaceAll(l["reader.email_body_send_code"], "{code}", validateCode))

	smtpPort, _ := strconv.Atoi(settings["app.smtpPort"])
	d := gomail.NewDialer(settings["app.smtpHost"], smtpPort, settings["app.smtpEmail"], settings["app.smtpPassword"])
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": l["reader.email_send_failed"]})
		return
	}

	components.Cache.Set(c.validateCodeKey(req.Email), validateCode, 10*time.Minute)

	ctx.JSON(http.StatusOK, rsp)
}

type ReaderApiLoginReq struct {
	Email        string `json:"email"`
	ValidateCode string `json:"validateCode"`
}

type ReaderApiLoginResult int

const (
	ReaderLoginResultOK         ReaderApiLoginResult = 0
	ReaderLoginResultInputError ReaderApiLoginResult = 1
	ReaderLoginResultWrongCode  ReaderApiLoginResult = 2
	ReaderLoginResultNewReader  ReaderApiLoginResult = 3
	ReaderLoginResultInactive   ReaderApiLoginResult = 4
)

type ReaderApiLoginRsp struct {
	Result  ReaderApiLoginResult `json:"result"`
	Message string               `json:"message"`
	Reader  *models.Reader       `json:"reader"`
}

func (c *ReaderApiController) GenerateToken(ctx *gin.Context, email string) {
	readerTokenKey := c.readerTokenKey(email)
	readerToken := utils.RandString(64)
	components.Cache.Set(readerTokenKey, readerToken, 30*24*time.Hour)
	ctx.SetCookie(ReaderTokenCookieName, fmt.Sprintf("%s %s", email, readerToken), 30*24*60*60, "/", "", false, true)
}

func (c *ReaderApiController) Login(ctx *gin.Context) {
	settings := models.GetAllSettingsMap()
	l := locales.Locales[settings["app.language"]]

	req, rsp := &ReaderApiLoginReq{}, &ReaderApiLoginRsp{}
	if err := ctx.BindJSON(&req); err != nil {
		rsp.Result, rsp.Message = ReaderLoginResultInputError, err.Error()
		ctx.JSON(http.StatusBadRequest, rsp)
		return
	}

	validateCodeKey := c.validateCodeKey(req.Email)
	validateCode := components.Cache.Get(validateCodeKey)
	if req.ValidateCode != validateCode {
		rsp.Result, rsp.Message = ReaderLoginResultWrongCode, l["reader.wrong_validate_code"]
		ctx.JSON(http.StatusBadRequest, rsp)
		return
	}

	reader := models.GetReaderByEmail(req.Email)
	if reader == nil {
		components.Cache.Set(validateCodeKey, validateCode, 24*time.Hour)
		rsp.Result = ReaderLoginResultNewReader
		ctx.JSON(http.StatusOK, rsp)
		return
	}

	components.Cache.Del(validateCodeKey)

	if !reader.IsActive {
		rsp.Result, rsp.Message = ReaderLoginResultInactive, l["reader.inactive"]
		ctx.JSON(http.StatusForbidden, rsp)
		return
	}

	c.GenerateToken(ctx, req.Email)
	rsp.Result = ReaderLoginResultOK
	ctx.JSON(http.StatusOK, rsp)
}

func (c *ReaderApiController) Logout(ctx *gin.Context) {
	if email, ok := ctx.Get(ReaderEmail); ok {
		readerTokenKey := c.readerTokenKey(email.(string))
		components.Cache.Del(readerTokenKey)
	}
	ctx.SetCookie(ReaderTokenCookieName, "", 0, "/", "", false, true)
}

type ReaderApiRegisterReq struct {
	Nickname     string     `json:"nickname" binding:"alphanumunicode"`
	Gender       uint8      `json:"gender"`
	Birthday     *time.Time `json:"birthday"`
	Email        string     `json:"email" binding:"email"`
	Phone        string     `json:"phone" binding:"numeric"`
	ValidateCode string     `json:"validateCode"`
}

type ReaderApiRegisterRsp struct {
	Message string `json:"message"`
}

func (c *ReaderApiController) Register(ctx *gin.Context) {
	settings := models.GetAllSettingsMap()
	l := locales.Locales[settings["app.language"]]

	var (
		req ReaderApiRegisterReq
		rsp ReaderApiRegisterRsp
	)
	if err := ctx.ShouldBindJSON(&req); err != nil {
		rsp.Message = err.Error()
		ctx.JSON(http.StatusBadRequest, rsp)
		return
	}

	validateCodeKey := c.validateCodeKey(req.Email)
	validateCode := components.Cache.Get(validateCodeKey)
	if req.ValidateCode != validateCode {
		rsp.Message = l["reader.wrong_validate_code"]
		ctx.JSON(http.StatusBadRequest, rsp)
		return
	}

	if models.GetReaderByEmail(req.Email) != nil {
		rsp.Message = l["reader.exists_email"]
		ctx.JSON(http.StatusBadRequest, rsp)
		return
	}
	if models.GetReaderByNickname(req.Nickname) != nil {
		rsp.Message = l["reader.exists_nickname"]
		ctx.JSON(http.StatusBadRequest, rsp)
	}

	components.Cache.Del(validateCodeKey)

	reader := &models.Reader{
		Nickname: req.Nickname,
		Gender:   req.Gender,
		Birthday: req.Birthday,
		Email:    req.Email,
		Phone:    req.Phone,
		IsActive: true,
	}
	reader.Save()

	c.GenerateToken(ctx, req.Email)
	ctx.JSON(http.StatusOK, rsp)
}

func (c *ReaderApiController) Info(ctx *gin.Context) {
	if r, ok := ctx.Get(ReaderKey); ok && r != nil {
		if reader := models.GetReaderByEmail(ctx.GetString(ReaderEmail)); reader.IsActive {
			ctx.JSON(http.StatusOK, reader)
			return
		}
	}
	ctx.JSON(http.StatusUnauthorized, gin.H{})
}

type ReaderApiInfoModifyReq struct {
	Nickname string     `json:"nickname" binding:"alphanumunicode"`
	Gender   uint8      `json:"gender"`
	Birthday *time.Time `json:"birthday"`
	Phone    string     `json:"phone" binding:"numeric"`
}

type ReaderApiInfoModifyRsp struct {
	Message string `json:"message"`
}

func (c *ReaderApiController) InfoModify(ctx *gin.Context) {
	var (
		req ReaderApiInfoModifyReq
		rsp ReaderApiInfoModifyRsp
	)
	if err := ctx.ShouldBindJSON(&req); err != nil {
		rsp.Message = err.Error()
		ctx.JSON(http.StatusBadRequest, rsp)
		return
	}

	reader := models.GetReaderByEmail(ctx.GetString(ReaderEmail))
	reader.Nickname = req.Nickname
	reader.Gender = req.Gender
	reader.Birthday = req.Birthday
	reader.Phone = req.Phone
	reader.Update()

	ctx.JSON(http.StatusOK, rsp)
}

func (c *ReaderApiController) Comment(ctx *gin.Context) {
	ctx.JSON(http.StatusNotImplemented, gin.H{})
}

func (c *ReaderApiController) Action(ctx *gin.Context) {
	ctx.JSON(http.StatusNotImplemented, gin.H{})
}

func NewReaderApiController(r gin.IRouter) *ReaderApiController {
	c := &ReaderApiController{}

	router := r.Group("/reader")
	router.Use(c.Middleware)
	router.POST("/code", c.SendCode)
	router.POST("/login", c.Login)
	router.POST("/logout", c.Logout)
	router.POST("/register", c.Register)
	router.GET("/info", c.Info)
	router.POST("/info", c.InfoModify)
	router.POST("/comment", c.Comment)
	router.POST("/action", c.Action)

	return c
}
