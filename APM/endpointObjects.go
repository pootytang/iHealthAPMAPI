package apm

type AllEPs struct {
	Endpoints []struct {
		FilesEndpoints    FilesEndpointList    `json:"files_endpoints,omitempty"`
		CommandsEndpoints CommandsEndpointList `json:"commands_endpoints,omitempty"`
	} `json:"endpoints"`
}

type FilesEndpointList []struct {
	Endpoint FilesEndPoint
}

type CommandsEndpointList []struct {
	Endpoint CommandsEndpoint
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

// Endpoint
type EndPoint struct {
	Eps []EPData `json:"endpoint"`
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
	Path        string `json:"path,omitempty"`
	FullPath    string `json:"full_path"`
	Description string `json:"description"`
}
