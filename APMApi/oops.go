package apmapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func Oops(ctx *gin.Context) {
	log.Warnf("Oops page accessed from client ip: %s", ctx.ClientIP())
	ctx.JSON(http.StatusNotFound, gin.H{"message": "Oops, you need a qkview id when accessing this endpoint"})
	return
}
