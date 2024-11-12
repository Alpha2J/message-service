package incoming_http

type CreateEmailReq struct {
	From    string   `json:"from" validate:"required,email"`
	To      []string `json:"to" validate:"required,dive,email"` // 验证每个元素都是有效的邮箱
	Subject string   `json:"subject" validate:"required"`
	Body    string   `json:"body" validate:"required"`
}

type CreateWechatWorkMessageReq struct {
	FromAppID int64    `json:"from_app_id" validate:"required,gt=0"` // 验证值需要大于0
	To        []string `json:"to" validate:"required"`
	Content   string   `json:"content" validate:"required"`
}
