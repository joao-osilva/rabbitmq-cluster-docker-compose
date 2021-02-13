#!/bin/sh

openssl req -x509 -nodes -newkey rsa:1024 -keyout  mykey.pem -out mykey.pem
