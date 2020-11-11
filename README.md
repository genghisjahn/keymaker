### Keymaker makes RSA keys and JSON Web Tokens (JWTs).

`keymaker` is a command line utility written in Go.  It has two functions:
1. Create RSA key pairs.
1. Create a JWT based on those key pairs.

The JWT could then be used to interact with service APIs.  The private key file is used to create the JWT and the public key (.pub) file is sent to the service APIs so that the JWT can be verified.  This way only the private key holder can create valid JWTs, and the public key can only be used to verify them.  This allows the two different parties to establish verification without sharing a private key.  Please see Miguel Greenberg's [blog post](https://blog.miguelgrinberg.com/post/json-web-tokens-with-public-key-signatures) and this helpful [Stack Overflow](https://stackoverflow.com/a/44352675/13324985) for more information.  If you like to dig into `mono-spaced` specification documents, try [this](https://tools.ietf.org/html/rfc7518#page-8) out.

`./keymaker -help` will display the various command line flags that can be passed to `keymaker`.

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

#### Creating a JWT using an existing RSA private key

You'll need to specify a few command line arguments to generate a valid JWT.  

1. `-keyfile` the name of the private RSA key that will be used to sign the token.  This is the file that does _not_ end in .pub.
1. `-sub` string value is the _subject_ of the JWT.  Typically the source/signer of the JWT.
1. `-aud` string value is the _audience_ of the JWT.  Typically the audience is the service that will be verifying the JWT before taking some action.
1. `-exp` integer value is the number of hours that will be added to the current time.  This future date will be the expiration date of the JWT.  So, `-exp 1` will create a token that will be valid for 1 hour before it expires.
1. `-jwtfile` string value is the name of the file that will containt the JWT token.  The suffic `.jwt` will be applied to this argument.

When creating a JWT, `keymaker` will add Unix time stamps for `iat` (Issue At) and `nbf` (Not Before).  The values will be the time the JWT was created.


The header of the JWT will look like this:
```json
{
  "alg": "RS256",
  "typ": "JWT"
}
```

The payload of the JWT will look something like this:
```json
{
  "aud": "temp_aud",
  "exp": 1605491395,
  "iat": 1605131395,
  "nbf": 1605131395,
  "sub": "temp_sub"
}
```

The JWT file itself won't look like that JSON output, it'll be a long line of base64 encoded parts separated by periods(.).  To see the values you'll have to decode it (https://jwt.io/).



#### Creating an RSA Key pair and a JWT in one step
Lastly, you can do both steps in one command.  That is, you can create an RSA key pair and create a jwt file.  For example:

`./keymaker -name temp_key  -aud temp_aud -sub temp_sub  -exp 100 -jwt temp_jwt` 

This command will create the following four files:

* `test_rsa.rsa` --This is the private key
* `test_rsa.pub` --This is the public key
* `test.rsa.pub.base64` --This is the public key encoded in base64
* `test_jwt.jwt` --This is the file that contains the JWT
