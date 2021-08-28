package middleware

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

type LogMiddlewareConfig struct {
	Path string
}

func NewLogMiddleware(c *LogMiddlewareConfig) gin.HandlerFunc {
	LOG.path = c.Path
	LOG.logDataQ = make(chan *logData, 10000)
	go LOG.proccessLogData()

	return func(ctx *gin.Context) {
		tm := time.Now().Local()
		ctx.Next()
		data := &logData{
			tm:        tm,
			ip:        ctx.ClientIP(),
			method:    ctx.Request.Method,
			uri:       ctx.Request.RequestURI,
			protocol:  ctx.Request.Proto,
			status:    ctx.Writer.Status(),
			size:      ctx.Writer.Size(),
			referer:   ctx.Request.Referer(),
			userAgent: ctx.Request.UserAgent(),
		}

		LOG.logDataQ <- data
	}
}

type logData struct {
	tm        time.Time
	ip        string
	method    string
	uri       string
	protocol  string
	status    int
	size      int
	referer   string
	userAgent string
}

type logNS struct {
	path     string
	file     *os.File
	writer   *bufio.Writer
	fileName string
	logDataQ chan *logData
}

var LOG logNS

func (ns *logNS) getWriter(fileName string) *bufio.Writer {
	if fileName == ns.fileName {
		return ns.writer
	}

	if ns.file != nil {
		ns.writer = nil
		ns.file.Close()
	}

	var err error
	filePath := filepath.Join(ns.path, fileName)
	ns.file, err = os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalln(err)
	}
	ns.writer = bufio.NewWriter(ns.file)

	return ns.writer
}

func (ns *logNS) proccessLogData() {
	for {
		data := <-ns.logDataQ

		ymd := data.tm.Format("2006-01-02")
		fileName := fmt.Sprintf("elf-%s.log", ymd)
		writer := ns.getWriter(fileName)
		date := data.tm.Format("2/Jan/2006:15:04:05 -0700")

		logMsg := fmt.Sprintf(`%s [%s] "%s %s %s" %d %d "%s" "%s"`,
			data.ip, date, data.method, data.uri, data.protocol, data.status,
			data.size, data.referer, data.userAgent)
		writer.WriteString(logMsg + "\n")
		writer.Flush()
	}
}
