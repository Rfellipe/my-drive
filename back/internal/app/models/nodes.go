package models

// How nodes are stored
type FsNode struct {
	ID         *string `json:"id"`
	Name       string  `json:"name" binding:"required"`
	Type       string  `json:"type" binding:"required"`
	Size       *int    `json:"size"`
	ParentId   *string `json:"parentId"`
	Created_at *string `json:"created_at"`
	Updated_at *string `json:"updated_at"`
}
