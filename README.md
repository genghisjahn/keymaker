## keymaker

### Makes RSA keys and JSON Web Tokens (JWTs).

Keymake is a command line utility written in Go.  It has two functions:
1. Create RSA key pairs.
1. Create a JWT based on those key pairs.

The JWT could then be used to interact with service APIs.  The private key file is used to create the JWT and the public key (.pub) file is sent to the service APIs so that the JWT can be verified.  This one only the private key holder can create valid JWTs, and the public key can only be used to verify them.  This allows the two different parties to establish verification without sharing private keys.  Please see Miguel Greenberg's [blog post](https://blog.miguelgrinberg.com/post/json-web-tokens-with-public-key-signatures) and this helpful [Stack Overflow](https://stackoverflow.com/a/44352675/13324985) for more information.  If you like to dig into mono-spaced specification documents, try [this](https://tools.ietf.org/html/rfc7518#page-8) out.

#### Generate a Keypair

Once the source has been compiled, run the executable like this:

`./keymaker -name test_rsa`

This will give you an output like this:

```text
2020/11/10 17:40:47 No JWT created.
2020/11/10 17:40:47 test_rsa.pub.base64
```

Three files are created:

* `test_rsa.rsa` --This is the private key
* `test_rsa.pub` --This is the public key
* `test.rsa.pub.base64` --This is the public key encoded in base64

Why the base64 encoding?
