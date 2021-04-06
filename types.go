package main

type JSONKeyInfo struct {
	PrivateKeyPath string `json:"private_key_path"`
	Subject        string `json:"sub"`
	Audience       string `json:"aud"`
	Issuer         string `json:"iss"`
	Scope          string `json:"scope"`
	Expiration     int    `json:"exp"`
	JWTFile        string `json:"jwt_file"`
}

/*
	name := flag.String("name", "temp", "The base name of the private key file to be used to sign the JWT. If the file is called private.rsa you would just enter 'private'.")
	keyfile := flag.String("keyfile", "temp_file", "The name of an existing private RSA key to use to sign a JWT.")
	bsize := flag.Int("size", 4096, "Bitsize of the RSA key.  The default is 4096.")
	sub := flag.String("sub", "", "Subject(sub) for the JWT.  If left blank no JWT will be created.  The subject is the kinds of services/data that will be acted upon of the JWT.")
	aud := flag.String("aud", "", "Audience(aud) for the JWT.  If left blank no JWT will be created.  This audience is the service that will be verifying and extracting data from the JWT to do something.")
	iss := flag.String("iss", "", "Issuer(iss) for the JWT.  If left blank no JWT will be created. This issue is the entity that creates the JWT.")
	scope := flag.String("scope", "", "Scope(scope) for the JWT.  If left blank no JWT will be created.  The scope is a space delimited value that dictates what the JWT can do.")
	exp := flag.Int("exp", 0, "Expiration(exp) hours from current unix time for the JWT expiration. If left blank no JWT will be created.")
	jwtfile := flag.String("jwt", "", "The name of file that will contain the jwt token.  The suffix '.jwt' will be appended to this value.  If left blank no JWT will be created.")
*/
