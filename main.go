package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"

	b64 "encoding/base64"
	"encoding/pem"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/ssh"
)

// openssl genrsa -out app.rsa keysize
// openssl rsa -in app.rsa -pubout > app.rsa.pub

func main() {
	name := flag.String("name", "temp", "The base name of the private key file to be used to sign the JWT. If the file is called private.rsa you would just enter 'private'.")
	bsize := flag.Int("size", 2048, "Bitsize of the RSA key.  The default is 4096.")
	jwtjson := flag.String("jwtjson", "", "The name of the json file that contains properties used to create the JWT file.")
	flag.Parse()

	j := JSONKeyInfo{}
	if *jwtjson != "" {

		loadErr := j.LoadFromFile(*jwtjson)
		if loadErr != nil {
			log.Fatalln(loadErr)

		}
	}
	if *name != "temp" {
		if err := makeRSAKeys(*name, *bsize); err != nil {
			fmt.Println(err)
			return
		}
		savePubKeyToBase64(*name)
		return
	}
	if errV := j.IsValid(); errV == nil {
		jwt, jErr := makeJWT(j)
		if jErr != nil {
			fmt.Println(jErr)
		}
		saveJWT(j.JWTFile, jwt)
	} else {
		fmt.Println("Invalid json file: ", errV.Error())
	}
}

func makeJWT(j JSONKeyInfo) (string, error) {
	n := time.Now()
	signBytes, err := ioutil.ReadFile(j.PrivateKeyPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	signKey, keyErr := ssh.ParseRawPrivateKey(signBytes)

	if keyErr != nil {
		fmt.Println(keyErr)
	}
	ctime := n.Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iss":   j.Issuer,
		"sub":   j.Subject,
		"nbf":   ctime,
		"exp":   n.Add(time.Hour * time.Duration(j.Expiration)).Unix(),
		"aud":   j.Audience,
		"scope": j.Scope,
		"iat":   ctime,
	})

	// Sign and get the complete encoded token as a string using the secret
	return token.SignedString(signKey)

}

func saveJWT(filename, token string) error {
	outFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer outFile.Close()
	_, err2 := io.WriteString(outFile, token)
	if err2 != nil {
		return err2
	}
	return outFile.Sync()
}

func savePubKeyToBase64(name string) error {
	filename := name + ".rsa.pub"
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	sEnc := b64.StdEncoding.EncodeToString(data)
	err2 := ioutil.WriteFile(name+".pub.base64", []byte(sEnc), 0644)
	return err2
}

func makeRSAKeys(filename string, size int) error {
	/*
		Golang RSA Code provide by https://stackoverflow.com/a/64105068/13324985
		Much thanks!
	*/

	// Generate RSA key.
	key, err := rsa.GenerateKey(rand.Reader, size)
	if err != nil {
		return err
	}

	// Extract public component.
	pub := key.Public()

	// Encode private key to PKCS#1 ASN.1 PEM.

	keyPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(key),
		},
	)

	// Encode public key to PKCS#1 ASN.1 PEM.

	pubBytes, pbErr := x509.MarshalPKIXPublicKey(pub.(*rsa.PublicKey))
	if pbErr != nil {
		log.Fatal(pbErr)
	}
	pubPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "PUBLIC KEY",
			Bytes: pubBytes,
		},
	)
	// Write private key to file.
	if err := ioutil.WriteFile(filename+".rsa", keyPEM, 0700); err != nil {
		return err
	}

	// Write public key to file.
	if err := ioutil.WriteFile(filename+".rsa.pub", pubPEM, 0755); err != nil {
		return err
	}
	return nil
}
