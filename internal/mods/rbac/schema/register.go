package schema

type RegisterRequest struct {
	Nickname string `json:"nickname" form:"nickname" binding:"required,min=1,max=32"`
	Password string `json:"password" form:"password" binding:"required,min=8,max=72"`
	Email    string `json:"email" form:"email" binding:"required,email"`
	Code     string `json:"code" form:"code" binding:"required,len=6,numeric"`
}

type RegisterID struct {
	ID int64 `json:"id"`
}

type RegisterResponse struct {
	Success bool       `json:"success"`
	Data    RegisterID `json:"data"`
}
