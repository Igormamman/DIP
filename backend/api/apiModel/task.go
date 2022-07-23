package apiModel

type TaskBatch struct {
	Data []Task `json:"data"`
}

type Task struct {
	Url   string `json:"url"`
	Token string `json:"tag"`
}

type BaseJson struct {
	Type string `json:"type"`
}

type ClientResultResponse struct {
	BaseJson
	Data []TaskResult `json:"data"`
}

type TaskResult struct {
	Token        string   `json:"image"`
	DetectedNums []string `json:"detected_nums"`
	DebugImageUrl *string  `json:"debug_image_url,omitempty"`
}

type NewFile struct {
	FilePath string `json:"filePath"`
	FileSize int    `json:"fileSize"`
	Token    string `json:"token"`
}

