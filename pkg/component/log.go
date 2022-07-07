package component

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/mingslife/bone"

	"elf-server/pkg/conf"
	"elf-server/pkg/utils"
)

type logData struct {
	tm       time.Time
	ip       string
	method   string
	uri      string
	protocol string
	// status    int
	// size      int
	referer   string
	userAgent string
}

type Log struct {
	Router   *bone.Router `inject:"application.router"`
	path     string
	file     *os.File
	writer   *bufio.Writer
	fileName string
	logDataQ chan *logData
}

func (*Log) Name() string {
	return "component.log"
}

func (*Log) Init() error {
	return nil
}

func (c *Log) Register() error {
	cfg := conf.GetConfig()
	if cfg.Log {
		c.path = LogPath()
		c.logDataQ = make(chan *logData, 1000)
		go c.proccessLogData()
		c.Router.Use(c.Middleware)
	}
	return nil
}

func (*Log) Unregister() error {
	return nil
}

func (c *Log) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			return
		}

		tm := time.Now().Local()
		next.ServeHTTP(w, r)
		data := &logData{
			tm:       tm,
			ip:       utils.ClientIP(r),
			method:   r.Method,
			uri:      r.RequestURI,
			protocol: r.Proto,
			// status:    -1,
			// size:      -1,
			referer:   r.Header.Get("Referer"),
			userAgent: r.Header.Get("User-Agent"),
		}
		c.logDataQ <- data
	})
}

func (c *Log) getWriter(fileName string) *bufio.Writer {
	if fileName == c.fileName {
		return c.writer
	}

	if c.file != nil {
		c.writer = nil
		c.file.Close()
	}

	var err error
	filePath := filepath.Join(c.path, fileName)
	c.file, err = os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalln(err)
	}
	c.writer = bufio.NewWriter(c.file)

	return c.writer
}

func (c *Log) proccessLogData() {
	for {
		data := <-c.logDataQ

		ymd := data.tm.Format("2006-01-02")
		fileName := LogFileName(ymd)
		writer := c.getWriter(fileName)
		date := data.tm.Format("2/Jan/2006:15:04:05 -0700")

		logMsg := fmt.Sprintf(`%s [%s] "%s %s %s" "%s" "%s"`,
			data.ip, date, data.method, data.uri, data.protocol,
			data.referer, data.userAgent)
		writer.WriteString(logMsg + "\n")
		writer.Flush()
	}
}

func LogPath() string {
	return path.Join("log")
}

func LogFileName(ymd string) string {
	return fmt.Sprintf("access-%s.log", ymd)
}

var _ bone.Component = (*Log)(nil)
