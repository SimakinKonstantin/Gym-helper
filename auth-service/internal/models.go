package internal

type SignInReq struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type SignInResp struct {
	Token string `json:"token"`
}

type SignUpReq struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UserDb struct {
	Login string `db:"login"`
	Hash  string `db:"hash"`
}
