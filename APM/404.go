package apm

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func Oops(ctx *gin.Context) {
	// setup some vars to shorten lines
	path := ctx.Request.URL.Path
	ip := ctx.ClientIP()
	url := ctx.Request.URL
	log.Warnf("Oops %s accessed from client ip %s for page %s is invalid", url, ip, path)

	// Check which path was accessed
	var msg string
	switch true {
	case strings.HasPrefix(path, "/apm/commands"):
		log.Warnf("/apm/commands endpoint accessed with wrong path: %s", path)
		msg = fmt.Sprintf("OOPS, invalid commands endpoint: %s. See %s for endpoints ", path, url.Hostname()+path)
	case strings.HasPrefix(path, "/apm/files"):
		log.Warnf("/apm/files endpoint accessed with wrong path: %s", path)
		msg = fmt.Sprintf("OOPS, invalid files endpoint: %s. See %s for endpoints ", path, url.Hostname()+path)
	default:
		log.Warn("Invalid Endpoint: %s", url)
		msg = fmt.Sprintf("OOPS, invalid commands endpoint: %s. See %s for all endpoints ", path, url.Hostname()+path)
	}

	ctx.JSON(http.StatusNotFound, gin.H{"code": "404", "message": msg})
}
