#!/bin/bash

#   How to use this shell script....
#
#
#   makejwt.sh uses 5 parameters
#
#
#   #1 is the kid or Key ID that identifies the key pair used to sign/identify the token.
#
#   $2 is the path to the file that contains the proto-jwt claims that will be put into the final, signed JWT
#      There are place holder values in this file for iat, nbf and exp
#   
#   $3 is the number of hours that will be added to the current time for the JWT's expiration
#
#   $4 is double quote enclosed (") space delimited values for scope "scope1 scope2 scope3"
#
#   $5 is the path to the private rsa key that will be used to sign the JWT
# 
#   $6 is name file that the signed JWT will be output to
#
#   Example:

#   ./climakejwt.sh 12345678 proto.json 10 "scope1 scope2" example.rsa example.jwt
#   
#   The above command will set the "kid" value in the header to 12345678,
#   pen the proto.json ($2) file, 
#   replace place holder values for iat, nbf with the current unix/epoch time,
#   replace place holder value exp  with current time + 10 ($3) hours as unix/epoch time,
#   put the header, payload and signtaure together into one string, 
#   create the signature based on that string with the private key in the file example.rsa ($5),
#   output the result to the file example.jwt ($6).

jwt_header=$(echo -n '{"alg":"RS256","typ":"JWT","kid":"'$1'"}' | base64 | sed s/\+/-/ | sed -E s/=+$//)
t=$(date +%s)
exptime=$(date -j -v +$3H +%s)
scope=$4

filepayload=`cat $2`

filepayload=${filepayload/\"iat_value\"/$t}
filepayload=${filepayload/\"nbf_value\"/$t}
filepayload=${filepayload/\"exp_value\"/ $exptime}
filepayload=${filepayload/scope_value/$scope}


payload=$(echo -n $filepayload | base64 | sed s/\+/-/ | sed -E s/=+$//)

body=${jwt_header}.${payload}

signature=$(echo -n $body | openssl dgst -sha256 -binary -sign $5  | openssl enc -base64 | tr -d '\n=' | tr -- '+/' '-_')
  
jwt=${jwt_header}.${payload}.${signature}
echo $jwt > $6