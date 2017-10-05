#!/bin/bash

tag=beekman9527/webdriver

#1. get webdriver_exporter
webdriver=webdriver.linux
cp ../_output/$webdriver ./webdriver.Dockerfile

#2. get promethus
prometheus=prometheus
mkdir $prometheus
rm ./$prometheus/*
wget https://github.com/prometheus/prometheus/releases/download/v1.7.2/prometheus-1.7.2.linux-amd64.tar.gz -O $prometheus.tar.gz
tar xzf ${prometheus}.tar.gz -C $prometheus --strip-components=1
cp ./prometheus.yml $prometheus

#3. get chromedriver
rm chromedriver_linux64.zip chromedriver_linux64
wget http://chromedriver.storage.googleapis.com/2.21/chromedriver_linux64.zip
unzip chromedriver_linux64.zip

#4. run Docker build
docker build -t $tag .
