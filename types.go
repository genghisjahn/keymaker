package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"strings"
)

type JSONKeyInfo struct {
	PrivateKeyPath string `json:"private_key_path"`
	Subject        string `json:"sub"`
	Audience       string `json:"aud"`
	Issuer         string `json:"iss"`
	Scope          string `json:"scope"`
	Expiration     int    `json:"exp"`
	JWTFile        string `json:"jwt_file"`
}

func (j *JSONKeyInfo) LoadFromFile(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	jErr := json.Unmarshal(data, &j)
	if jErr != nil {
		return jErr
	}
	return nil
}

func (j JSONKeyInfo) IsValid() error {
	if len(j.PrivateKeyPath) == 0 {
		return errors.New("private_key_path value is missing")
	}
	if len(j.Subject) == 0 {
		return errors.New("sub value is missing")
	}
	if len(j.Audience) == 0 {
		return errors.New("aud value is missing")
	}
	if len(j.Issuer) == 0 {
		return errors.New("iss value is missing")
	}
	if len(j.Scope) == 0 {
		return errors.New("scope value is missing")
	}
	if j.Expiration <= 0 {
		return errors.New("exp must be greater than 0")
	}
	if len(j.JWTFile) == 0 {
		return errors.New("jwt_file value is missing")
	}
	if !strings.HasSuffix(j.JWTFile, ".jwt") {
		return errors.New("jwt_file value must end with .jwt")
	}
	return nil
}
