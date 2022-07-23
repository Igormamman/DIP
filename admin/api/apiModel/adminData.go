package apiModel

import (
	dbModel "admin/models"
)

type UrlParam struct {
	Table  string
	Label  string
}

type AdminHtmlPageObject struct {
	UrlParams []UrlParam
	Photo     []dbModel.Photo
}
