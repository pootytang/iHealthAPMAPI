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
		Path:        "/id",
		FullPath:    "/apm/commands/id",
		Description: "returns all the apm related commands in qkview id",
	}
	data = append(data, d)

	//ID with slash
	d = EPData{
		Method:      "GET",
		Params:      "null",
		Path:        "/id/",
		FullPath:    "/apm/commands/<id>/",
		Description: "returns all the apm related commands in qkview id",
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
		Path:        "/id/command_id",
		FullPath:    "/apm/command/id/command_id",
		Description: "returns the result of running command_id from qkview id",
	}
	data = append(data, d)

	//ID with slash
	d = EPData{
		Method:      "GET",
		Params:      "null",
		Path:        "/id/command_id/",
		FullPath:    "/apm/command/id/command_id/",
		Description: "returns the result of running command_id from qkview id",
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
		Path:        "/id",
		FullPath:    "/apm/files/id",
		Description: "returns all the files endpoints",
	}
	data = append(data, d)

	//ID with slash
	d = EPData{
		Method:      "GET",
		Params:      "null",
		Path:        "/id/",
		FullPath:    "/apm/files/id/",
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
	var commandsEndpoint CommandsEndpoint
	var commandEndpoint CommandEndpoint
	var filesEndpoint FilesEndPoint
	var fileEndpoint FileEndPoint
	var authEndpoint AuthEndPoint
	var metaDataEndpoint MetadataEndPoint
	var baseEndpoint BaseEndPoint

	/********** COMMANDS **********/
	var commandsData []EPData
	// Empty path
	cmds := EPData{
		Method:      "GET",
		Params:      "null",
		Path:        "",
		FullPath:    "/apm/commands",
		Description: "returns all the commands endpoints",
	}
	commandsData = append(commandsData, cmds)

	//ID no slash
	cmds = EPData{
		Method:      "GET",
		Params:      "null",
		Path:        "/id",
		FullPath:    "/apm/commands/id",
		Description: "returns all the apm related commands in qkview id",
	}
	commandsData = append(commandsData, cmds)
	commandsEndpoint.Eps = commandsData

	/********** COMMAND **********/
	var commandData []EPData
	// Empty path
	cmd := EPData{
		Method:      "GET",
		Params:      "null",
		Path:        "",
		FullPath:    "/apm/command",
		Description: "returns all the command endpoints",
	}
	commandData = append(commandData, cmd)

	//ID no slash
	cmd = EPData{
		Method:      "GET",
		Params:      "null",
		Path:        "/id/command_id",
		FullPath:    "/apm/command/id/command_id",
		Description: "returns the result of running command_id from qkview id",
	}
	commandData = append(commandData, cmd)
	commandEndpoint.Eps = commandData

	/********** FILES **********/
	var filesData []EPData
	// Empty path
	fs := EPData{
		Method:      "GET",
		Params:      "null",
		Path:        "",
		FullPath:    "/apm/files",
		Description: "returns all the files endpoints",
	}
	filesData = append(filesData, fs)

	//ID no slash
	fs = EPData{
		Method:      "GET",
		Params:      "null",
		Path:        "/id",
		FullPath:    "/apm/files/id",
		Description: "returns all the files endpoints",
	}
	filesData = append(filesData, fs)
	filesEndpoint.Eps = filesData

	/********** FILE **********/
	var fileData []EPData
	// Empty path
	f := EPData{
		Method:      "GET",
		Params:      "null",
		Path:        "",
		FullPath:    "/apm/file",
		Description: "returns the file endpoints",
	}
	fileData = append(fileData, f)

	// POST Empty path
	f = EPData{
		Method:      "POST",
		Params:      `{"id":"qkview_id", "file_id":"file_id", "file_name":"name_of_file"}`,
		Path:        "",
		FullPath:    "/apm/file",
		Description: "Post body contains json params. Returns the requested file",
	}
	fileData = append(fileData, f)
	fileEndpoint.Eps = fileData

	/********** AUTH **********/
	var authData []EPData
	auth := EPData{
		Method:      "POST",
		Params:      `{"user_id":"f5 username", "user_secret":"password"}`,
		Path:        "",
		FullPath:    "/auth",
		Description: "Post body contains json params. Returns json containing status and expiration",
	}
	authData = append(authData, auth)
	authEndpoint.Eps = authData

	/********** METADATA **********/
	var metaData []EPData
	// Empty path
	md := EPData{
		Method:      "GET",
		Params:      "id",
		Path:        "/id",
		FullPath:    "/metadata/id",
		Description: "Returns the metadata for the given qkview id",
	}
	metaData = append(metaData, md)
	metaDataEndpoint.Eps = metaData

	/********** BASE **********/
	var baseData []EPData
	// Empty path
	bd := EPData{
		Method:      "GET",
		Params:      "",
		Path:        "",
		FullPath:    "/",
		Description: "Returns all the metadata for the app",
	}
	baseData = append(baseData, bd)
	baseEndpoint.Eps = baseData

	//Prepare the all endpoints objest
	ep := Endpoints{
		Commands: commandsEndpoint,
		Command:  commandEndpoint,
		Files:    filesEndpoint,
		File:     fileEndpoint,
		Auth:     authEndpoint,
		Metadata: metaDataEndpoint,
		Base:     baseEndpoint,
	}

	var eps []Endpoints
	eps = append(eps, ep)
	all := AllEndPoints{Endpoints: eps}
	// all := AllEndPoints{
	// 	{Endpoints: ep},
	// }
	ctx.JSON(http.StatusOK, all)
}
