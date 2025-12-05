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
	SessionId    string `json:"sessionId,omitempty"`
	CsrfToken    string `json:"csrfToken,omitempty"`
	// RefreshToken 仅用于在后端写 HttpOnly Cookie，不会返回给前端
	RefreshToken string `json:"-"`
}

type ListUserResp struct {
	List []UserResp `json:"list"`
}

// 刷新 Token 请求无需 Body，依赖 HttpOnly refreshToken Cookie
type RefreshTokenReq struct{}

type RefreshTokenResp struct {
	AccessToken  string `json:"accessToken"`
	AccessExpire int64  `json:"accessExpire"`
	Id           string `json:"id"`
	Username     string `json:"username"`
	Role         string `json:"role"`
	SessionId    string `json:"sessionId,omitempty"`
	DeviceId     string `json:"deviceId,omitempty"`
}

type VerifyTokenResp struct {
	Valid bool      `json:"valid"`
	User  *UserResp `json:"user,omitempty"`
}
