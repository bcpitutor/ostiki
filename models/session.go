package models

type Session struct {
	SessID       string `json:"sessId"`
	SessionOwner string `json:"sessionOwner"`
	AccessToken  string `json:"accessToken"`
	IdToken      string `json:"idToken"`
	Expire       string `jspn:"expire"`
	RefreshToken string `json:"refreshToken"`
	TokenType    string `json:"tokenType"`
	UserInfo     struct {
		Id            string `json:"id"`
		Email         string `json:"email"`
		VerifiedEmail string `json:"verifiedEmail"`
		Picture       string `json:"picture"`
		Hd            string `json:"hd"`
	}
	Epoch           int64  `json:"epoch"`
	SessionExpEpoch int64  `json:"sessionExpEpoch"`
	PreviousIdToken string `json:"prevIdToken"`
	Details         string `json:"details"`
	Rtimes          int    `json:"rTimes,default:1"`
	IsRevoked       bool   `json:"isRevoked,default:false"`
}

type SessionExpose struct {
	SessionId       string `json:"sessionId"`
	SessionOwner    string `json:"sessionOwner"`   //  email
	SessionDetails  string `json:"sessionDetails"` // comment
	ExpiresAt       string `json:"expiresAt"`
	Epoch           int64  `json:"epoch"`
	SessionExpEpoch int64  `json:"sessionExpEpoch"`
	RefreshCount    int    `json:"refreshCount"`
	Revoked         bool   `json:"revoked"`
}
