package types

type CreateUserReq struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResp struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

type DeleteUserReq struct {
	Id string `path:"id"`
}

type LoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResp struct {
	AccessToken  string `json:"accessToken"`
	AccessExpire int64  `json:"accessExpire"`
	Id           string `json:"id"`
	Username     string `json:"username"`
	Role         string `json:"role"`
}

type ListUserResp struct {
	List []UserResp `json:"list"`
}
