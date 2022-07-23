package apiModel

import (
	"mime/multipart"
)

type Photo struct {
	// to properly read form in stream part sequence is important
	ImageSize              int             `form:"image_size"`
	ResizedSize            int             `form:"resized_size"`
	MLSize                 int             `form:"ml_size"`
	WatermarkSize          int             `form:"watermark_size"`
	ImageWidth             int             `form:"image_width"`
	ImageHeight            int             `form:"image_height"`
	UID                    string          `form:"uid"`
	ImageUpload            *multipart.Part `form:"image"`
	ResizedImageUpload     *multipart.Part `form:"resized_image"`
	WatermarkedImageUpload *multipart.Part `form:"watermarked_image"`
	ResizedMLImageUpload   *multipart.Part `form:"resizedML_image"`
}

type PhotoMeta struct {
	Height      int      `json:"height"`
	Width       int      `json:"width"`
	Competitors string   `json:"competitors"`
	Race        RaceMeta `json:"Race"`
}

type RaceMeta struct {
	Name string `json:"name"`
	Date uint64 `json:"date"`
	City string `json:"city"`
}
