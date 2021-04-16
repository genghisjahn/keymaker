#!/bin/bash

./keymaker  -keyfile  files/dcuniverse.rsa -sub userlink -aud dcuniverse -iss dcuniverse -scope userlink-* -exp 1000 -jwt dcuniverse
./keymaker  -keyfile  files/boomerang.rsa -sub userlink -aud boomerang -iss boomerang -scope userlink-* -exp 1000 -jwt boomerang
./keymaker  -keyfile  files/kiosk.rsa -sub location -aud wiwzorld -iss sorting-hat-kiosk -scope onsite-process -exp 1000 -jwt kiosk
./keymaker  -keyfile  files/lotr.rsa -sub userlink -aud lotr -iss lotr -scope userlink-* -exp 1000 -jwt lotr
./keymaker  -keyfile  files/matrix.rsa -sub userlink -aud matrix -iss matrix -scope userlink-* -exp 1000 -jwt matrix
./keymaker  -keyfile  files/wbgames-batmanlego.rsa -sub userlink -aud wbgames -iss wbg-batmanlego -scope userlink-* -exp 1000 -jwt wbbatmanlego
./keymaker  -keyfile  files/wbgames-hl.rsa -sub userlink -aud wbgames -iss wbg-hl -scope userlink-* -exp 1000 -jwt wbg-hl
./keymaker  -keyfile  files/wbgames-hplego.rsa -sub userlink -aud wbgames -iss wbg-hplego -scope userlink-* -exp 1000 -jwt wbg-hplego
./keymaker  -keyfile  files/wizworld.rsa -sub userlink -aud wizworld -iss wizworld -scope userlink-* onsite-create -exp 1000 -jwt wizworld