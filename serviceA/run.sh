#!/bin/bash

docker build --build-arg DEBUG_FLAG="-DDEBUG" -t servicea .
docker save -o servicea.tar servicea
gzip servicea.tar
