package main

import (
	"github.com/gin-gonic/gin"
	apm "github.com/pootytang/iHealthAPMAPI/APM"
	apmapi "github.com/pootytang/iHealthAPMAPI/APMApi"
	ihealthapi "github.com/pootytang/iHealthAPMAPI/iHealthAPI"
	log "github.com/sirupsen/logrus"
)

// ihealthApiMidware configures the api client as a gin Context middleware
// this is so the cookies are re-used between requests to the APM Api Frontend
func ihealthApiMidware(c *ihealthapi.ApiClient) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set("apiclient", c)
		ctx.Next()
	}
}

func main() {
	log.SetLevel(log.DebugLevel)
	c := ihealthapi.NewClient()
	router := gin.New()
	router.Use(ihealthApiMidware(c))

	// NOTE:  addomg the equivalent / routes avoids gin responding with a 307
	// THE BASE ROUTES
	router.GET("/", apm.AllEndpoints)
	router.POST("/", apm.AllEndpoints)
	router.POST("/auth", apmapi.Auth)
	router.POST("/auth/", apmapi.Auth)
	router.GET("/metadata", apmapi.Oops)
	router.GET("/metadata/:id", apmapi.Metadata)
	router.GET("/endpoints", apm.AllEndpoints)
	router.GET("/endpoints/", apm.AllEndpoints)

	// APM ROUTES
	apmCommandsGroup := router.Group("/apm/commands")
	{
		apmCommandsGroup.GET("", apm.CommandsEndpoints)
		apmCommandsGroup.GET("/", apm.CommandsEndpoints)
		apmCommandsGroup.GET("/:id", apm.Commands)
		apmCommandsGroup.GET("/:id/", apm.Commands)
	}

	apmCommandGroup := router.Group("/apm/command")
	{
		apmCommandGroup.GET("", apm.CommandEndpoints)
		apmCommandGroup.GET("/", apm.CommandEndpoints)
		apmCommandGroup.GET("/:id", apm.Oops)
		apmCommandGroup.GET("/:id/", apm.Oops)
		apmCommandGroup.GET("/:id/:commandId", apm.Command)
		apmCommandGroup.GET("/:id/:commandId/", apm.Command)
	}

	apmFilesGroup := router.Group("/apm/files")
	{
		apmFilesGroup.GET("", apm.FilesEndpoints)
		apmFilesGroup.GET("/", apm.FilesEndpoints)
		apmFilesGroup.GET("/:id", apm.Files)
		apmFilesGroup.GET("/:id/", apm.Files)
	}

	apmFileGroup := router.Group("/apm/file")
	{
		apmFileGroup.GET("", apm.FileEndpoints)
		apmFileGroup.GET("/", apm.FileEndpoints)
		apmFileGroup.POST("", apm.File)
		apmFileGroup.POST("/", apm.File)
	}

	apmProfilesGroup := router.Group("/apm/profiles")
	{
		apmProfilesGroup.GET("", apm.Profiles)
		apmProfilesGroup.GET("/", apm.Profiles)
		apmProfilesGroup.POST("", apm.Profiles)
		apmProfilesGroup.POST("/", apm.Profiles)
	}
	router.Run(":8080")
}
