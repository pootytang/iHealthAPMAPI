package apm

// This code works but doesn't look how I want it to
// type AllEndPoints []struct {
// 	Endpoints Endpoints `json:"endpoints"`
// }
// This was at the end of endpointMethods.go
// all := AllEndPoints{
// 	{Endpoints: ep},
// }

type AllEndPoints struct {
	Endpoints []Endpoints
}

//all := AllEndPoints{Endpoints: ep}

type Endpoints struct {
	Files    FilesEndPoint
	File     FileEndPoint
	Commands CommandsEndpoint
	Command  CommandEndpoint
	Auth     AuthEndPoint
	Metadata MetadataEndPoint
	Base     BaseEndPoint
}

/******** APM ********/
//Commands Plural
type CommandsEndpoint struct {
	Eps []EPData `json:"commands_endpoints"`
}

//Command Singular
type CommandEndpoint struct {
	Eps []EPData `json:"command_endpoints"`
}

//Files Plurral
type FilesEndPoint struct {
	Eps []EPData `json:"files_endpoints"`
}

//File Singular
type FileEndPoint struct {
	Eps []EPData `json:"file_endpoints"`
}

// Profiles Plural
type ProfilesEndPoint struct {
	Eps []EPData `json:"profiles_endpoints"`
}

/******************************/

/******** BASE ********/
/* these will be apart of the AllEndpoints method */
// Auth
type AuthEndPoint struct {
	Eps []EPData `json:"auth_endpoints"`
}

// Metadata
type MetadataEndPoint struct {
	Eps []EPData `json:"metadata_endpoints"`
}

// Base
type BaseEndPoint struct {
	Eps []EPData `json:"base_endpoint"`
}

/******************************/

// EPData is the base structure for all Endpoints
type EPData struct {
	Method      string `json:"method"`
	Params      string `json:"params"`
	Path        string `json:"path"`
	FullPath    string `json:"full_path"`
	Description string `json:"description"`
}
