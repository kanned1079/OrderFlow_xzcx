package dto

type CommitCommentByOrderIdRequestDto struct {
	OrderId     string   `json:"order_id"`
	UserId      int64    `json:"user_id"`
	CommentText string   `json:"comment_text"`
	ImagesUrls  []string `json:"images_urls"`
}

type FetchCommentListByMIdRequestDto struct {
	Page int `form:"page" json:"page"`
	Size int `form:"size" json:"size"`
}
