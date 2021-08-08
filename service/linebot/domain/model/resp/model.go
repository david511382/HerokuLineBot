package resp

type OAuthVerifyIDToken struct {
	// URL used to generate the ID token.
	Iss string `json:"iss"`
	// User ID for which the ID token was generated.
	Sub string `json:"sub"`
	// Channel ID
	Aud string `json:"aud"`
	// The expiry date of the ID token in UNIX time.
	Exp uint `json:"exp"`
	// Time when the ID token was generated in UNIX time.
	Lat uint `json:"iat"`
	// Time the user was authenticated in UNIX time. Not included if the max_age value wasn't specified in the authorization request.
	AuthTime uint `json:"auth_time"`
	// The nonce value specified in the authorization URL. Not included if the nonce value wasn't specified in the authorization request.
	Nonce string `json:"nonce"`
	// A list of authentication methods used by the user. One or more of:
	// pwd: Log in with email and password
	// lineautologin: LINE automatic login (including through LINE SDK)
	// lineqr: Log in with QR code
	// linesso: Log in with single sign-on
	Amr []string `json:"amr"`
	// User's display name. Not included if the profile scope wasn't specified in the authorization request.
	Name string `json:"name"`
	// User's profile image URL. Not included if the profile scope wasn't specified in the authorization request.
	Picture string `json:"picture"`
	// User's email address. Not included if the email scope wasn't specified in the authorization request.
	Email string `json:"email"`
}

type GetUserProfile struct {
	DisplayName string `json:"displayName"`
}

type ReplyMessage struct{}

type PushMessage struct{}

type ListRichMenuRichMenu struct {
	Name       string `json:"name"`
	RichMenuID string `json:"richMenuId"`
}

type ListRichMenu struct {
	RichMenus []*ListRichMenuRichMenu `json:"richmenus"`
}

type DeleteRichMenu struct{}

type SetDefaultRichMenu struct{}

type SetRichMenuTo struct{}

type SetRichMenuTos struct{}

type GetDefaultRichMenu struct {
	RichMenuID string `json:"richMenuId"`
}

type UploadRichMenuImage struct{}

type CreateRichMenu struct {
	RichMenuID string `json:"richMenuId"`
}
