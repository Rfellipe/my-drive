package models

// PUT payload for route under /dir/:id
type DirUpdate struct {
	Node    FsNode `json:"node"`
	OldPath string `json:"oldPath"`
	NewPath string `json:"newPath"`
}

// DELETE payload for route under /dir/:id
type DirDel struct {
	Node FsNode `json:"node"`
	Soft bool   `json:"soft"`
}
