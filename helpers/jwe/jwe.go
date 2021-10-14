package jwe

import (
	"encoding/json"
	"github.com/lestrrat/go-jwx/jwa"
	"github.com/lestrrat/go-jwx/jwe"
)

// Credential ...
type Credential struct {
	KeyLocation string
	Passphrase  string
}

// Generate ...
func (cred *Credential) GenerateJwePayload(id string) (res string, err error) {
	privkey, err := rsaConfigSetup(cred.KeyLocation, cred.Passphrase)
	if err != nil {
		return res, err
	}

	// Generate payload
	payload := map[string]interface{}{
		"id": id,
	}
	payloadString, err := json.Marshal(payload)
	if err != nil {
		return res, err
	}

	// Generate JWE
	jweRes, err := jwe.Encrypt([]byte(payloadString), jwa.RSA1_5, &privkey.PublicKey, jwa.A128CBC_HS256, jwa.Deflate)
	res = string(jweRes)

	return res, err
}

// Rollback ...
func (cred *Credential) Rollback(id string) (res map[string]interface{}, err error) {
	privkey, err := rsaConfigSetup(cred.KeyLocation, cred.Passphrase)
	if err != nil {
		println(err.Error())
		return res, err
	}

	decrypted, err := jwe.Decrypt([]byte(id), jwa.RSA1_5, privkey)
	if err != nil {
		return res, err
	}

	res = map[string]interface{}{}
	err = json.Unmarshal(decrypted, &res)
	if err != nil {
		return res, err
	}

	return res, err
}
