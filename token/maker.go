package token

type Maker interface {
	CreateToken()
	VerifyToken()
}
