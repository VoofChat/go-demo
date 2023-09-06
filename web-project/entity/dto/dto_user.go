package dto

type UserRegister struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

type UserQuery struct {
	UserId uint64 `json:"userId" form:"userId"`
}

type UserUpdate struct {
	UserId   uint64 `json:"userId" form:"userId"`
	Username string `json:"username" form:"username"`
}
