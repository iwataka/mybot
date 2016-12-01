#!/bin/bash

openssl genrsa 2048 > mybot.key
openssl req -new -key mybot.key > mybot.csr
openssl x509 -days 3650 -req -signkey mybot.key < mybot.csr > mybot.crt
