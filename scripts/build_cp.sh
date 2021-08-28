#!/bin/bash


echo -e "\033[0;32m[build & cp  kraken-testfs]\033[0m"
go build -o kraken-testfs ./../tools/bin/testfs
sudo cp ./kraken-testfs  /usr/local/bin/

echo -e "\033[0;32m[build & cp  kraken-proxy]\033[0m"
go build -o kraken-push-cli ./../push-cli
sudo cp ./kraken-push-cli  /usr/local/bin/

echo -e "\033[0;32m[build & cp  kraken-origin]\033[0m"
go build -o kraken-origin ./../origin
sudo cp ./kraken-origin  /usr/local/bin/

echo -e "\033[0;32m[build & cp  kraken-tracker]\033[0m"
go build -o kraken-tracker ./../tracker
sudo cp ./kraken-tracker  /usr/local/bin/

echo -e "\033[0;32m[build & cp  kraken-build-index]\033[0m"
go build -o kraken-build-index ./../build-index
sudo cp ./kraken-build-index  /usr/local/bin/

