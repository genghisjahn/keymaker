package main

import (
	"crypto/x509"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func main() {
	name := flag.String("name", "temp", "The of the private key file to be used to sign the JWT")
	aud := flag.String("aud", "none", "Audeince(aud) for the JWT.  If left blank no JWT will be created.  This is typcally the service that will be verifying and extracting data from the JWT to do something.")
	sub := flag.String("sub", "none", "Subject(sub) for the JWT.  If left blank no JWT will be created.  This is typcally the calling service that is supplying the JWT to the specified AUD service.")
	exp := flag.Int("exp", 0, "Expiration(exp) hours from current unix time for the JWT expiration. If left blank no JWT will be created.")
	flag.Parse()

	privName := *name + "rsa"

	if len(*aud) > 0 && len(*sub) > 0 && *exp > 0 {
		jwt, jErr := makeJWT(privName, *aud, *sub, *exp)
		if jErr != nil {
			log.Println(jErr)
		}
		fmt.Println(jwt)
	}

}

func makeJWT(privepath, aud, sub string, exp int) (string, error) {
	n := time.Now()
	signBytes, err := ioutil.ReadFile(privepath)
	if err != nil {
		log.Fatal(err)
	}
	signKey, keyErr := x509.ParsePKCS1PrivateKey(signBytes)
	if keyErr != nil {
		log.Println(keyErr)
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"sub": sub,
		"nbf": n.Unix(),
		"exp": n.Add(time.Hour * time.Duration(exp)),
		"aud": aud,
		"iat": n.Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	return token.SignedString(signKey)

}

func saveJWT(filename, token string) {
	outFile, err := os.Create(filename)
	checkError(err)
	defer outFile.Close()
	_, err2 := io.WriteString(outFile, token)
	checkError(err2)
	outFile.Sync()
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}
