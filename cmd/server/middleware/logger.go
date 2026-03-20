package middleware

import (
	"os"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

func InitLogger() {
	os.MkdirAll("logs", 0755)
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(&lumberjack.Logger{
		Filename:   "./logs/requests.log",
		MaxSize:    5,
		MaxBackups: 3,
		MaxAge:     14,
		Compress:   true,
	})
}

func JSONLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.WithFields(log.Fields{
			"ip":     c.ClientIP(),
			"method": c.Request.Method,
			"path":   c.Request.URL.Path,
		}).Info("Incoming Request")
		c.Next()
	}
}
