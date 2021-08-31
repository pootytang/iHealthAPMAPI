package apm

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// Return all commands related endpoints
func CommandsEndpoints(ctx *gin.Context) {
	log.Info("CommandsEndpoint() -> called")
	var data []EPData
	var commandsEndpoint CommandsEndpoint

	// Empty path
	d := EPData{
		Method:      "GET",
		Params:      "null",
		Path:        "",
		FullPath:    "/apm/commands",
		Description: "returns all the commands endpoints",
	}
	data = append(data, d)

	// Slash
	d = EPData{
		Method:      "GET",
		Params:      "null",
		Path:        "/",
		FullPath:    "/apm/commands/",
		Description: "returns all the commands endpoints",
	}
	data = append(data, d)

	//ID no slash
	d = EPData{
		Method:      "GET",
		Params:      "null",
		Path:        "/<id>",
		FullPath:    "/apm/commands/<id>",
		Description: "returns all the apm related commands in qkview <id>",
	}
	data = append(data, d)

	//ID with slash
	d = EPData{
		Method:      "GET",
		Params:      "null",
		Path:        "/<id>/",
		FullPath:    "/apm/commands/<id>/",
		Description: "returns all the apm related commands in qkview <id>",
	}
	data = append(data, d)

	commandsEndpoint.Eps = data
	ctx.JSON(http.StatusOK, commandsEndpoint)
}

// Return all commands related endpoints
func CommandEndpoints(ctx *gin.Context) {
	log.Info("CommandEndpoint() -> called")
	var data []EPData
	var commandEndpoint CommandEndpoint

	// Empty path
	d := EPData{
		Method:      "GET",
		Params:      "null",
		Path:        "",
		FullPath:    "/apm/command",
		Description: "returns all the command endpoints",
	}
	data = append(data, d)

	// Slash
	d = EPData{
		Method:      "GET",
		Params:      "null",
		Path:        "/",
		FullPath:    "/apm/command/",
		Description: "returns all the command endpoints",
	}
	data = append(data, d)

	//ID no slash
	d = EPData{
		Method:      "GET",
		Params:      "null",
		Path:        "/<id>/<command_id>",
		FullPath:    "/apm/command/<id>/<command_id>",
		Description: "returns the result of running <command_id> from qkview <id>",
	}
	data = append(data, d)

	//ID with slash
	d = EPData{
		Method:      "GET",
		Params:      "null",
		Path:        "/<id>/<command_id>/",
		FullPath:    "/apm/command/<id>/<command_id>/",
		Description: "returns the result of running <command_id> from qkview <id>",
	}
	data = append(data, d)

	commandEndpoint.Eps = data
	ctx.JSON(http.StatusOK, commandEndpoint)
}

// Return all files related endpoints
func FilesEndpoints(ctx *gin.Context) {
	log.Info("FilesEndpoint() -> called")
	var data []EPData
	var filesEndpoint FilesEndPoint

	// Empty path
	d := EPData{
		Method:      "GET",
		Params:      "null",
		Path:        "",
		FullPath:    "/apm/files",
		Description: "returns all the files endpoints",
	}
	data = append(data, d)

	// Slash
	d = EPData{
		Method:      "GET",
		Params:      "null",
		Path:        "/",
		FullPath:    "/apm/files",
		Description: "returns all the files endpoints",
	}
	data = append(data, d)

	//ID no slash
	d = EPData{
		Method:      "GET",
		Params:      "null",
		Path:        "/<id>",
		FullPath:    "/apm/files/<id>",
		Description: "returns all the files endpoints",
	}
	data = append(data, d)

	//ID with slash
	d = EPData{
		Method:      "GET",
		Params:      "null",
		Path:        "/<id>/",
		FullPath:    "/apm/files/<id>/",
		Description: "returns all the files endpoints",
	}
	data = append(data, d)

	filesEndpoint.Eps = data
	ctx.JSON(http.StatusOK, filesEndpoint)
}

// Return all file related endpoints
func FileEndpoints(ctx *gin.Context) {
	log.Info("FileEndpoint() -> called")
	var data []EPData
	var fileEndpoint FileEndPoint

	// Empty path
	d := EPData{
		Method:      "GET",
		Params:      "null",
		Path:        "",
		FullPath:    "/apm/file",
		Description: "returns the file endpoints",
	}
	data = append(data, d)

	// With slash
	d = EPData{
		Method:      "GET",
		Params:      "null",
		Path:        "/",
		FullPath:    "/apm/files/",
		Description: "returns all the file endpoints",
	}
	data = append(data, d)

	// POST Empty path
	d = EPData{
		Method:      "POST",
		Params:      `{"id":"qkview_id", "file_id":"file_id", "file_name":"name_of_file"}`,
		Path:        "",
		FullPath:    "/apm/file",
		Description: "Post body contains json params. Returns the requested file",
	}
	data = append(data, d)

	// POST With slash
	d = EPData{
		Method:      "POST",
		Params:      `{"id":"qkview_id", "file_id":"file_id", "file_name":"name_of_file"}`,
		Path:        "/",
		FullPath:    "/apm/file/",
		Description: "Post body contains json params. Returns the requested file",
	}
	data = append(data, d)

	fileEndpoint.Eps = data
	ctx.JSON(http.StatusOK, fileEndpoint)
}

// Return all endpoints
func ProfilesEndpoint(ctx *gin.Context) {

}

// Return all endpoints
func AllEndpoints(ctx *gin.Context) {

}
