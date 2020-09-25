openssl genrsa -out $1.rsa $2
openssl rsa -in $1.rsa -pubout > $1.rsa.pub