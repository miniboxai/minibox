package provider

import (
	"log"
	"testing"
	"time"

	osin "github.com/RangelReale/osin"
)

func TestJwtToken(t *testing.T) {
	var (
		tokenStr string
		err      error
		ad       osin.AccessData
	)

	tokenGen, _ := loadTokenGenJWT()
	tokenGen.SetClientCallback(func(cid string, data *osin.AccessData) error {
		data.Client = &osin.DefaultClient{
			Id:          "1234",
			Secret:      "asdfasdf",
			RedirectUri: "http://localhost:14000/auth",
		}
		return nil
	})

	client := &osin.DefaultClient{
		Id:          "1234",
		Secret:      "asdfasdf",
		RedirectUri: "http://localhost:14000/auth",
	}

	// 加密成 jwt
	if tokenStr, err = tokenGen.Encrypto(&osin.AccessData{
		Client:       client,
		AccessToken:  "1324123asdfasdf",
		CreatedAt:    time.Now(),
		ExpiresIn:    3600,
		RefreshToken: "1234dfasdferr12347daf",
	}); err != nil {
		t.Fatalf("Encrypto failed %s", err)
	} else {
		log.Printf("crypto string %s", tokenStr)
	}

	// 加密 jwt
	if err = tokenGen.Decrypto(tokenStr, &ad); err != nil {
		t.Fatalf("Decrypto failed %s", err)
	} else {
		log.Printf("decrypt osin.AccessData string %#v", &ad)
	}
}
