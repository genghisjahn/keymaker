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
	"strings"
	"time"

	b64 "encoding/base64"
	"encoding/pem"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/ssh"
)

//openssl genrsa -out app.rsa keysize
// openssl rsa -in app.rsa -pubout > app.rsa.pub

func main() {
	name := flag.String("name", "temp", "The base name of the private key file to be used to sign the JWT. If the file is called private.rsa you would just enter 'private'.")
	//	keyfile := flag.String("keyfile", "temp_file", "The name of an existing private RSA key to use to sign a JWT.")
	bsize := flag.Int("size", 4096, "Bitsize of the RSA key.  The default is 4096.")
	//	sub := flag.String("sub", "", "Subject(sub) for the JWT.  If left blank no JWT will be created.  The subject is the kinds of services/data that will be acted upon of the JWT.")
	//	aud := flag.String("aud", "", "Audience(aud) for the JWT.  If left blank no JWT will be created.  This audience is the service that will be verifying and extracting data from the JWT to do something.")
	//	iss := flag.String("iss", "", "Issuer(iss) for the JWT.  If left blank no JWT will be created. This issue is the entity that creates the JWT.")
	//	scope := flag.String("scope", "", "Scope(scope) for the JWT.  If left blank no JWT will be created.  The scope is a space delimited value that dictates what the JWT can do.")
	//	exp := flag.Int("exp", 0, "Expiration(exp) hours from current unix time for the JWT expiration. If left blank no JWT will be created.")
	//	jwtfile := flag.String("jwt", "", "The name of file that will contain the jwt token.  The suffix '.jwt' will be appended to this value.  If left blank no JWT will be created.")
	jwtfile := flag.String("jwtfile", "", "The name of the json file that contains properties used to create the JWT file.")
	flag.Parse()
	privkeyname := ""

	if *name != "temp" {
		privkeyname = *name + ".rsa"
	}

	lowJWT := strings.ToLower(*jwtfile)
	if strings.HasSuffix(lowJWT, ".jwt") {
		lowJWT = strings.TrimSuffix(lowJWT, ".jwt")
		jwtfile = &lowJWT
	}

	j := JSONKeyInfo{}

	if *name != "temp" {

		if err := makeRSAKeys(*name, *bsize); err != nil {
			fmt.Println(err)
			return
		}
		savePubKeyToBase64(*name)
	}
	if len(j.Audience) > 0 && j.Expiration > 0 && len(j.Subject) > 0 && len(j.JWTFile) > 0 && len(j.Issuer) > 0 {
		jwt, jErr := makeJWT(privkeyname, j.Issuer, j.Audience, j.Subject, j.Scope, j.Expiration)
		if jErr != nil {
			fmt.Println(jErr)
		}
		saveJWT(*jwtfile, jwt)
	} else {
		fmt.Println("Invalid json file")
	}
}

func makeJWT(privepath, iss, aud, sub, scope string, exp int) (string, error) {
	n := time.Now()
	signBytes, err := ioutil.ReadFile(privepath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	signKey, keyErr := ssh.ParseRawPrivateKey(signBytes)

	if keyErr != nil {
		fmt.Println(keyErr)
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iss":   iss,
		"sub":   sub,
		"nbf":   n.Unix(),
		"exp":   n.Add(time.Hour * time.Duration(exp)).Unix(),
		"aud":   aud,
		"scope": scope,
		"iat":   n.Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	return token.SignedString(signKey)

}

func saveJWT(filename, token string) {
	outFile, err := os.Create(filename + ".jwt")
	checkError(err)
	defer outFile.Close()
	_, err2 := io.WriteString(outFile, token)
	checkError(err2)
	outFile.Sync()
}

func savePubKeyToBase64(name string) {
	filename := name + ".rsa.pub"
	data, err := ioutil.ReadFile(filename)
	checkError(err)
	sEnc := b64.StdEncoding.EncodeToString(data)
	err2 := ioutil.WriteFile(name+".pub.base64", []byte(sEnc), 0644)
	checkError(err2)
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
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
