package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"

	b64 "encoding/base64"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/ssh"
)

//openssl genrsa -out app.rsa keysize
// openssl rsa -in app.rsa -pubout > app.rsa.pub

func main() {
	name := flag.String("name", "temp", "The of the private key file to be used to sign the JWT")
	aud := flag.String("aud", "none", "Audeince(aud) for the JWT.  If left blank no JWT will be created.  This is typcally the service that will be verifying and extracting data from the JWT to do something.")
	sub := flag.String("sub", "none", "Subject(sub) for the JWT.  If left blank no JWT will be created.  This is typcally the calling service that is supplying the JWT to the specified AUD service.")
	exp := flag.Int("exp", 0, "Expiration(exp) hours from current unix time for the JWT expiration. If left blank no JWT will be created.")
	flag.Parse()

	privName := *name + ".rsa"

	if len(*aud) > 0 && len(*sub) > 0 && *exp > 0 {
		jwt, jErr := makeJWT(privName, *aud, *sub, *exp)
		if jErr != nil {
			log.Println(jErr)
		}
		saveJWT(*name, jwt)
	}
	savePubKeyToBase64(*name)

}

func makeJWT(privepath, aud, sub string, exp int) (string, error) {
	n := time.Now()
	signBytes, err := ioutil.ReadFile(privepath)
	if err != nil {
		log.Fatal(err)
	}
	signKey, keyErr := ssh.ParseRawPrivateKey(signBytes)

	if keyErr != nil {
		log.Println(keyErr)
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
