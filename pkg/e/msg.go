package e

var MsgFlags = map[int]string{
	SUCCESS:                         "ok",
	ERROR:                           "fail",
	INVALID_PARAMS:                  "Invalid parameters",
	ERROR_EXIST_TAG:                 "Tag name already exists",
	ERROR_EXIST_TAG_FAIL:            "Failed to get existing tag",
	ERROR_NOT_EXIST_TAG:             "Tag does not exist",
	ERROR_GET_TAGS_FAIL:             "Failed to get all tags",
	ERROR_COUNT_TAG_FAIL:            "Failed to count tags",
	ERROR_ADD_TAG_FAIL:              "Failed to add tag",
	ERROR_EDIT_TAG_FAIL:             "Failed to modify tag",
	ERROR_DELETE_TAG_FAIL:           "Failed to delete tag",
	ERROR_EXPORT_TAG_FAIL:           "Failed to export tag",
	ERROR_IMPORT_TAG_FAIL:           "Failed to import tag",
	ERROR_NOT_EXIST_ARTICLE:         "Article does not exist",
	ERROR_ADD_ARTICLE_FAIL:          "Failed to add article",
	ERROR_DELETE_ARTICLE_FAIL:       "Failed to delete article",
	ERROR_CHECK_EXIST_ARTICLE_FAIL:  "Failed to check if article exists",
	ERROR_EDIT_ARTICLE_FAIL:         "Failed to modify article",
	ERROR_COUNT_ARTICLE_FAIL:        "Failed to count articles",
	ERROR_GET_ARTICLES_FAIL:         "Failed to get multiple articles",
	ERROR_GET_ARTICLE_FAIL:          "Failed to get article",
	ERROR_GEN_ARTICLE_POSTER_FAIL:   "Failed to generate article poster",
	ERROR_AUTH_CHECK_TOKEN_FAIL:     "Token authentication failed",
	ERROR_AUTH_CHECK_TOKEN_TIMEOUT:  "Token has expired",
	ERROR_AUTH_TOKEN:                "Failed to generate token",
	ERROR_AUTH:                      "Invalid username or password",
	ERROR_UPLOAD_SAVE_IMAGE_FAIL:    "Failed to save image",
	ERROR_UPLOAD_CHECK_IMAGE_FAIL:   "Failed to check image",
	ERROR_UPLOAD_CHECK_IMAGE_FORMAT: "Image validation error, problem with format or size",
}

// GetMsg get error information based on Code
func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}

	return MsgFlags[ERROR]
}
