#!/bin/bash

# This bash file takes 2 arguments:
#
# $1 is the name of the private key file that will be created
#   The public key will be create at $1.pub
#   The base64 encoded version of the public key will be create at $1.pub.base64

# $2 is the number specifying the keysize.  It should be at least 2048.
#
#
# The private key file is used to sign things, and the public key file is used to verify that the 
# signature was indeed signed by the private key.
#
# The private key should be kept secret.
# The public key can be shared.

# Example:
#
#   ./makekeys.sh keyfiles/test1.rsa 2048
#
#   The makekeys.sh shell script will create a private key at the location keyfiles/test1.rsa
#   It will have a keysize of 2048 bits
#   The public key will be created at the location keyfiles/test1.rsa.pub
#   The base64 encoded public key will be created at keyfiles/test1.rsa.pub.base64
#   A kid or Key ID

openssl genrsa -out $1 $2
openssl rsa -in $1 -pubout > $1.pub
base64 $1.pub > $1.pub.base64
kid=$(xxd -l 4 -c 4 -p < /dev/random)
echo "You must put this value in the kid (key id) of the JWT: " $kid

