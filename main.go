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
	keyfile := flag.String("keyfile", "temp_file", "The name of an existing private RSA key to use to sign a JWT.")
	bsize := flag.Int("size", 4096, "Bitsize of the RSA key.  The default is 4096.")
	sub := flag.String("sub", "", "Subject(sub) for the JWT.  If left blank no JWT will be created.  The subject is typically the source/signer of the JWT.")
	aud := flag.String("aud", "", "Audience(aud) for the JWT.  If left blank no JWT will be created.  This is typcally the service that will be verifying and extracting data from the JWT to do something.")
	exp := flag.Int("exp", 0, "Expiration(exp) hours from current unix time for the JWT expiration. If left blank no JWT will be created.")
	jwtfile := flag.String("jwt", "", "The name of file that will contain the jwt token.  The suffix '.jwt' will be appended to this value.  If left blank no JWT will be created.")
	flag.Parse()
	privkeyname := ""
	if *name != "temp" && *keyfile != "temp_file" {
		fmt.Println("Only the -name flag or the -keyfile flag can be specified at the same time.")
		return
	}
	if *name != "temp" {
		privkeyname = *name + ".rsa"
	}
	if *keyfile != "temp_file" {
		privkeyname = *keyfile
	}
	lowJWT := strings.ToLower(*jwtfile)
	if strings.HasSuffix(lowJWT, ".jwt") {
		lowJWT = strings.TrimSuffix(lowJWT, ".jwt")
		jwtfile = &lowJWT
	}

	if *name != "temp" {

		if err := makeRSAKeys(*name, *bsize); err != nil {
			fmt.Println(err)
			return
		}
		savePubKeyToBase64(*name)
	}
	if len(*aud) > 0 && *exp > 0 && len(*sub) > 0 && len(*jwtfile) > 0 {
		jwt, jErr := makeJWT(privkeyname, *aud, *sub, *exp)
		if jErr != nil {
			fmt.Println(jErr)
		}
		saveJWT(*jwtfile, jwt)
	}
}

func makeJWT(privepath, aud, sub string, exp int) (string, error) {
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
		"sub": sub,
		"nbf": n.Unix(),
		"exp": n.Add(time.Hour * time.Duration(exp)).Unix(),
		"aud": aud,
		"iat": n.Unix(),
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
