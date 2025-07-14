#!/bin/bash

curl -v -X POST http://188.225.24.208:8000/registration \
-H "Content-type: application/json" \
-d '{
"firstname": "vasya",
"secondname": "pupkin",
"lastname": "pupkevich",
"email": "vasya@mail.ru",
"password": "pwd",
"role": "student",
"group": "ips"
}'
