package apiModel

import "mime/multipart"

type OutPhotoPreviewData struct {
	Count      int64    `json:"count"`
	PreviewURL []string `json:"previewURL"`
}

type LoadStruct struct {
	// to properly read form in stream part sequence is important
	FilePath string          `form:"filepath" json:"filepath"`
	FileType string          `form:"file_type" json:"file_type"`
	File     *multipart.Part `form:"file" json:"file"`
}

type UploadMeta struct {
	// to properly read form in stream part sequence is important
	PhotoName     string `json:"PhotoName"`
	RaceID        string
	UserID        string
	ImageSize     int `json:"imageSize"`
	ResizedSize   int `json:"resizedSize"`
	MLSize        int `json:"mlSize"`
	WatermarkSize int `json:"watermarkSize"`
	ImageWidth    int `json:"imageWidth"`
	ImageHeight   int `json:"imageHeight"`
}

type TaskData struct {
	FileData []FileData `json:"fileData"`
	RaceID   string     `json:"RaceID"`
}

type FileData struct {
	FileName string `json:"fileName"`
	UserID   string `json:"UserID,omitempty"`
	SID      string `json:"SID,omitempty"`
	RaceID   string `json:"RaceID,omitempty"`
}

type PathToUpload struct {
	Token    string `json:"token"`
	FileType string `json:"fileType"`
}

type PhotoMeta struct {
	UUID        string   `json:"UUID"`
	RUID        string   `json:"RUID"`
	Height      int      `json:"height"`
	Width       int      `json:"width"`
	Competitors []string `json:"competitors"`
	PrevUUID    []string `json:"prevUUID"`
	NextUUID    []string `json:"nextUUID"`
}
