package core

import (
	"crypto/ecdsa"
	"errors"
	"github.com/getlantern/deepcopy"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"time"
)

type VP struct {
	Context []string `json:"@context"`
	Id      string   `json:"id,omitempty"`
	Type    []string `json:"type,omitempty"`
	Holder  string   `json:"holder,omitempty"`

	// jwt의 token형식으로 저장한다.
	VerifiableCredential []string `json:"verifiableCredential"`
	Proof                *Proof   `json:"proof,omitempty"`

	Token string
}

// JWT를 위한 claim
type JwtClaimsForVP struct {
	jwt.StandardClaims

	Nonce string
	Vp    VP `json:"vp,omitempty"`
}

func NewVP(id string, typ []string, holder string, vcTokens []string) (*VP, error) {
	newVP := &VP{
		Context: []string{
			"https://www.w3.org/2018/credentials/v1",
		},
		Id:                   id,
		Type:                 typ,
		Holder:               holder,
		VerifiableCredential: vcTokens,
	}
	return newVP, nil
}

func (vp *VP) GenerateJWT(verificationId string, pvKey *ecdsa.PrivateKey) string {
	aud := ""
	exp := time.Now().Add(time.Minute * 5).Unix() //만료 시간. 현재 + 5분
	jti := uuid.NewString()                       // JWT ID
	iat := time.Now().Unix()
	nbf := iat
	iss := vp.Holder
	sub := "Verifiable Presentation"

	// Proof를 제거하고 JWT를 만들기 위해 복제한다.
	vpTmp := new(VP)
	deepcopy.Copy(vpTmp, vp)
	vpTmp.Proof = nil

	jwtClaims := JwtClaimsForVP{
		jwt.StandardClaims{
			Audience:  aud,
			ExpiresAt: exp,
			Id:        jti,
			IssuedAt:  iat,
			Issuer:    iss,
			NotBefore: nbf,
			Subject:   sub,
		},
		"qwasd!234",
		*vpTmp,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwtClaims)
	token.Header["kid"] = verificationId

	tokenString, err := token.SignedString(pvKey)

	if err != nil {

	}

	vp.Token = tokenString

	return tokenString
}

func (vp *VP) isVerify() (bool, error) {
	if vp.Token == "" {
		return false, errors.New("Token is empty.")
	}

	return true, nil
}
