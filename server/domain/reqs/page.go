package reqs

type Page struct {
	PageSize uint `json:"page_size" form:"page_size" binding:"required" url:"page_size"`
	// 1 開始
	PageIndex uint `json:"page_index" form:"page_index" binding:"required" url:"page_index"`
}
