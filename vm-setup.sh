#!/bin/bash

sudo apt install make git -y || sudo yum install make git -y;

curl -s -L https://golang.org/dl/go1.15.linux-amd64.tar.gz -O;
sudo tar -C /usr/local -xzf go1.15.linux-amd64.tar.gz;

echo 'export PATH=$PATH:/usr/local/go/bin' | sudo tee -a /etc/profile
