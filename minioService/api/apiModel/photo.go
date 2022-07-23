package apiModel

import (
	"mime/multipart"
)

type LoadStruct struct {
	// to properly read form in stream part sequence is important
	FilePath string          `form:"filepath" json:"filepath"`
	FileSize int             `form:"file_size" json:"file_size"`
	FileType string          `form:"file_type" json:"file_type"`
	File     *multipart.Part `form:"file" json:"file"`
}

type GetStruct struct {
	FilePath string `form:"filepath"`
}

type NewFile struct {
	FilePath string `json:"filePath"`
	FileName string `json:"fileName"`
	FileSize int    `json:"fileSize"`
	Token    string `json:"token"`
}
