package lib

type FsTypes int

const (
	Directory FsTypes = iota
	File
)

var FsTypeName = map[FsTypes]string{
	Directory: "dir",
	File:      "file",
}

func (ss FsTypes) String() string {
	return FsTypeName[ss]
}

type FileInfo struct {
	Name       string `json:"name"`
	Type       string `json:"type"`
	Size       int    `json:"size"`
	Created_at string `json:"created_at"`
	Updated_at string `json:"updated_at"`
}
