#!/bin/sh
set -e

# fetch dependencies
go get

# compile binary
go build

# install binary
if ! [ -d /usr/local/bin ]; then
    sudo mkdir -p /usr/local/bin
    sudo chmod 755 /usr/local/bin
fi
sudo cp mac2mqtt /usr/local/bin/
sudo chmod 755 /usr/local/bin/mac2mqtt

# install mqtt credentials file
sudo mkdir -p /etc/mac2mqtt
sudo cp mac2mqtt.yaml /etc/mac2mqtt/
sudo chmod -R 600 /etc/mac2mqtt

# configure to run at startup
sudo cp com.bessarabov.mac2mqtt.plist /Library/LaunchDaemons/
sudo launchctl load /Library/LaunchDaemons/com.bessarabov.mac2mqtt.plist
