package provider

import (
	"crypto/rsa"
	"fmt"

	osin "github.com/RangelReale/osin"
	jwt "github.com/dgrijalva/jwt-go"
)

type ClientCallback func(cid string, data *osin.AccessData) error

type AccessTokenGenJWT struct {
	PrivateKey     *rsa.PrivateKey
	PublicKey      *rsa.PublicKey
	clientCallback ClientCallback
}

func (c *AccessTokenGenJWT) GenerateAccessToken(data *osin.AccessData, generaterefresh bool) (accesstoken string, refreshtoken string, err error) {
	// generate JWT access token

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"cid": data.Client.GetId(),
		"exp": data.ExpireAt().Unix(),
	})

	accesstoken, err = token.SignedString(c.PrivateKey)
	if err != nil {
		return "", "", err
	}

	if !generaterefresh {
		return
	}

	// generate JWT refresh token
	token = jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"cid": data.Client.GetId(),
	})

	refreshtoken, err = token.SignedString(c.PrivateKey)
	if err != nil {
		return "", "", err
	}
	return
}

func (c *AccessTokenGenJWT) SetClientCallback(handle ClientCallback) {
	c.clientCallback = handle
}

func (c *AccessTokenGenJWT) Encrypto(data *osin.AccessData) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"cid": data.Client.GetId(),
		"tok": data.AccessToken,
		"rtk": data.RefreshToken,
		"exp": data.ExpireAt().Unix(),
	})

	jwtdata, err := token.SignedString(c.PrivateKey)
	if err != nil {
		return "", err
	}

	return jwtdata, nil
}

func (c *AccessTokenGenJWT) EncryptoMap(data map[string]interface{}) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"cid": data["client_id"],
		"tok": data["access_token"],
		"rtk": data["refresh_token"],
		"exp": data["expires_in"],
	})

	jwtdata, err := token.SignedString(c.PrivateKey)
	if err != nil {
		return "", err
	}

	return jwtdata, nil
}

func (c *AccessTokenGenJWT) Decrypto(tokenString string, data *osin.AccessData) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// claims, _ := token.Claims.(jwt.MapClaims)
		return c.PublicKey, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if c.clientCallback != nil {
			c.clientCallback(claims["cid"].(string), data)
		}

		data.AccessToken = claims["tok"].(string)
		data.RefreshToken = claims["rtk"].(string)
		data.ExpiresIn, _ = claims["exp"].(int32)
		return nil
	} else {
		return err
	}
}

func (c *AccessTokenGenJWT) DecryptoMap(tokenString string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// claims, _ := token.Claims.(jwt.MapClaims)
		return c.PublicKey, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
		// if c.clientCallback != nil {
		// 	c.clientCallback(claims["cid"].(string), data)
		// }

		// data.AccessToken = claims["tok"].(string)
		// data.RefreshToken = claims["rtk"].(string)
		// data.ExpiresIn, _ = claims["exp"].(int32)
		// return nil
	} else {
		return nil, err
	}
}

func loadTokenGenJWT() (*AccessTokenGenJWT, error) {
	var (
		accessTokenGenJWT AccessTokenGenJWT
		err               error
	)

	if accessTokenGenJWT.PrivateKey, err = jwt.ParseRSAPrivateKeyFromPEM(privatekeyPEM); err != nil {
		fmt.Printf("ERROR: %s\n", err)
		return nil, err
	}

	if accessTokenGenJWT.PublicKey, err = jwt.ParseRSAPublicKeyFromPEM(publickeyPEM); err != nil {
		fmt.Printf("ERROR: %s\n", err)
		return nil, err
	}

	return &accessTokenGenJWT, nil
}

func PublicKeyJWT() (*AccessTokenGenJWT, error) {
	var (
		accessTokenGenJWT AccessTokenGenJWT
		err               error
	)
	if accessTokenGenJWT.PublicKey, err = jwt.ParseRSAPublicKeyFromPEM(publickeyPEM); err != nil {
		fmt.Printf("ERROR: %s\n", err)
		return nil, err
	}

	return &accessTokenGenJWT, nil
}

var (
	privatekeyPEM = []byte("")
	publickeyPEM  = []byte(`
-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA+8YezW/5I62fDkw67vdm
VpAXVePI0/g7kxGWvaVknxn+T2YdkpmcmlD5WRGHd37cGGAW1hTOxlcOzpteuU31
ymKekWqQygYzllMUVkHr+d2F0oeODpz4QvnnXPS//1oUmXyJ2+oXf5VTQ4qr16ND
lTkuZwe56ALvk/c25ydN6SVq6e6Z3BllgZ+JlwlezLu3kWkWaM1Z3tJMr1T/frv2
diZF5EzfNm/9FL8llV2rhdDZu7OL1QMhuSiQfUvYlRLLNS3v1wdl+V4Wy7bWnJvk
XTiGzp4Ra5gA25XxAaYn8GC96SZUAJ8/2J2k+1swFNVk5OW6T0LsSB9WF50t1Zi/
uQIDAQAB
-----END PUBLIC KEY-----`)
)
