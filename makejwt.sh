#!/bin/bash

#   How to use this shell script....
#
#
#   makejwt.sh uses 4 parameters
#
#   $1 is the path to the file that contains the proto-jwt claims that will be put into the final, signed JWT
#      There are place holder values in this file for iat, nbf and exp
#   
#   $2 is the number of hours that will be added to the current time for the JWT's expiration
#
#   $3 is the path to the private rsa key that will be used to sign the JWT
# 
#   $4 is name file that the signed JWT will be output to
#
#   Example:

#   ./climakejwt.sh proto.json 10 example.rsa example.jwt
#   
#   The above command will open the proto.json ($1) file, 
#   replace place holder values for iat, nbf with the current unix/epoch time
#   replace place holder value exp  with current time + 10 ($2) hours as unix/epoch time
#   put the header, payload and signtaure together into one string, 
#   create the signature based on that string with the private key in the file example.rsa ($3)
#   output the result to the file example.jwt ($4)

jwt_header=$(echo -n '{"alg":"RS256","typ":"JWT"}' | base64 | sed s/\+/-/ | sed -E s/=+$//)
t=$(date +%s)
exptime=$(date -j -v +$2H +%s)
filepayload=`cat $1`
filepayload=${filepayload/\"iat_value\"/$t}
filepayload=${filepayload/\"nbf_value\"/$t}
filepayload=${filepayload/\"exp_value\"/ $exptime}
payload=$(echo -n $filepayload | base64 | sed s/\+/-/ | sed -E s/=+$//)

body=${jwt_header}.${payload}

signature=$(echo -n $body | openssl dgst -sha256 -binary -sign $3  | openssl enc -base64 | tr -d '\n=' | tr -- '+/' '-_')
  
jwt=${jwt_header}.${payload}.${signature}
echo $jwt > $4
