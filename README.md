## keymaker

### Makes RSA keys and JSON Web Tokens (JWTs).

Keymake is a command line utility written in Go.  It has two functions:
1. Create RSA key pairs.
1. Create a JWT based on those key pairs.

The JWT could then be used to interact with service APIs.  The private key file is used to create the JWT and the public key (.pub) file is sent to the service APIs so that the JWT can be verified.  This one only the private key holder can create valid JWTs, and the public key can only be used to verify them.  This allows the two different parties to establish verification without sharing private keys.  Please see Miguel Greenberg's [blog post](https://blog.miguelgrinberg.com/post/json-web-tokens-with-public-key-signatures) and this helpful [Stack Overflow](https://stackoverflow.com/a/44352675/13324985) for more information.  If you like to dig into mono-spaced specification documents, try [this](https://tools.ietf.org/html/rfc7518#page-8) out.

#### Generate a Keypair

Once the source has been compiled, run the executable like this:

`./keymaker -name test_rsa`

Three files are created:

* `test_rsa.rsa` --This is the private key
* `test_rsa.pub` --This is the public key
* `test.rsa.pub.base64` --This is the public key encoded in base64

There is an optional parameter `-size` where you can specify the size of the RSA key.  The default is `4096`.

Why the base64 encoding?

Quick aside here.  I wanted to use the public keys was in the environment variables for AWS Lambda functions.  But when you paste the environment variables into the AWS Console, the carriage returns are removed and the public key no longer parses correctly.  Then I decided to base64 encode the public key so that I could paste the value into the environment variable and not have to worry about missing carriage returns.  But _then_ I ran into a limitiation where you can only have 4000 bytes of data _in total_ for all environment variables in a given Lambda.  The base64 encoding of the public key is stored elsewhere.  It's not needed, but it solved a problem I was having at the time.

###Creating a JWT using an existing RSA private key

You'll need to specify a few command line arguments to generate a valid JWT.  

1. `-sub` This string value is the _subject_ of the JWT.  Typically the source/signer of the JWT.
1. `-aud` This string value is the _audience_ of the JWT.  Typically the audience is the service that will be verifying the JWT before taking some action.
1. `-exp` This integer value is the number of hours that will be added to the current time.  This future date will be the expiration date of the JWT.  So, `-exp 1` will create a token that will be valid for 1 hour before it expires.
1. `-jwtfile` This string value is the name of the file that will containt the JWT token.  The suffic `.jwt` will be applied to this argument.

```go
name := flag.String("name", "temp", "The of the base name of the private key file to be used to sign the JWT. If the file is called private.rsa you would just enter 'private'.")
	keyfile := flag.String("keyfile", "temp_file", "The name of an existing private RSA key to use to sign a JWT.")
	bsize := flag.Int("size", 4096, "Bitsize of the RSA key.  The default is 4096.")
	sub := flag.String("sub", "", "Subject(sub) for the JWT.  If left blank no JWT will be created.  The subject is typically the source/signer of the JWT.")
	aud := flag.String("aud", "", "Audience(aud) for the JWT.  If left blank no JWT will be created.  This is typcally the service that will be verifying and extracting data from the JWT to do something.")
	exp := flag.Int("exp", 0, "Expiration(exp) hours from current unix time for the JWT expiration. If left blank no JWT will be created.")
    jwtfile := flag.String("jwt", "", "The name of file that will contain the jwt token.  The suffix '.jwt' will be appended to this value.  If left blank no JWT will be created.")
```