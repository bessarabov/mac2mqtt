#!/bin/sh
set -e

# remove configuration to run at startup
sudo launchctl unload /Library/LaunchDaemons/com.bessarabov.mac2mqtt.plist
sudo rm /Library/LaunchDaemons/com.bessarabov.mac2mqtt.plist

# remove mqtt credentials file
sudo rm -r /etc/mac2mqtt

# remove binary
sudo rm /usr/local/bin/mac2mqtt

# maybe remove /usr/local/bin if we created it
if [ -z "$(ls -A /usr/local/bin)" ]; then
    sudo rmdir /usr/local/bin
fi
