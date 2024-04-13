# mac2mqtt

`mac2mqtt` is a program that allows viewing and controlling some aspects of computers running macOS via MQTT.

It publishes to MQTT:

 * current volume
 * volume mute state
 * battery charge percent
 * if macOS is connected to MQTT

You can send commands to MQTT to:

 * change volume
 * mute/unmute
 * put computer to sleep
 * shutdown computer
 * turn off display

## Overview

To control macOS via MQTT using this project you need several elements:

 * File `mac2mqtt` with compiled program (also known as 'binary' or 'executable')
 * File with configuration `mac2mqtt.yaml`. You need to write this file yourself, but you can use file `mac2mqtt.yaml` stored in this repository as a starting point
 * A system that automatically launches mac2mqtt after macOS restarts. You can run `mac2mqtt` manually without such system, but in such case you need to start `mac2mqtt` every time after restarting macOS. Manually running mac2mqtt is a good starting point when setting up the project.
 * MQTT server (it is often called MQTT Broker)
 * Some system what will read data from MQTT and send command to MQTT. Originally this project was created to make it possible to control macOS computer via [Home Assistant](https://www.home-assistant.io/), but any software that work with MQTT can be used with `mac2mqtt`.

It is recommended to put executable file and configuration file in your home directory in the subdirectory `mac2mqtt`:

```
/Users/USERNAME/mac2mqtt/
├── mac2mqtt
└── mac2mqtt.yaml
```

## Installation

There are several ways you can get `mac2mqtt` binary.

### Use pre-compiled binaries

Download the file from GitHub Releases section.

Make sure to download the correct file for your architecture:

 * `mac2mqtt_VERSION_arm64` — use this file for Apple Silicon Macs
 * `mac2mqtt_VERSION_x86_64` — use this file for Intel-based Macs

You need to make downloaded file executable `chmod +x FILE_NAME`.

### Compile the source code yourself

You need golang to be installed on your system — https://go.dev/doc/install

 1. Clone this repo
 2. `cd mac2mqtt/`
 3. `go build .`

This will create file `mac2mqtt` that you can run.

## Running

To run this program you need 2 files in a directory:

```
/Users/USERNAME/mac2mqtt/
├── mac2mqtt
└── mac2mqtt.yaml
```

Take `mac2mqtt.yaml` that is stored in this repository, but edit it, and put your data into it.

Then run `./mac2mqtt` in the terminal. You should see see somthing like this:

    $ ./mac2mqtt
    2021/04/12 10:37:28 Started
    2021/04/12 10:37:29 Connected to MQTT
    2021/04/12 10:37:29 Sending 'true' to topic: mac2mqtt/bessarabov-osx/status/alive

## Running in the background

You need `mac2mqtt.yaml` and `mac2mqtt` to be placed in the directory `/Users/USERNAME/mac2mqtt/`,
then you need to create file `/Library/LaunchDaemons/com.bessarabov.mac2mqtt.plist`:

```xml
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
    <dict>
        <key>Label</key>
        <string>com.bessarabov.mac2mqtt</string>
        <key>Program</key>
        <string>/Users/USERNAME/mac2mqtt/mac2mqtt</string>
        <key>WorkingDirectory</key>
        <string>/Users/USERNAME/mac2mqtt/</string>
        <key>RunAtLoad</key>
        <true/>
        <key>KeepAlive</key>
        <true/>
    </dict>
</plist>
```

And run:

    launchctl load /Library/LaunchDaemons/com.bessarabov.mac2mqtt.plist

(To stop you need to run `launchctl unload /Library/LaunchDaemons/com.bessarabov.mac2mqtt.plist`)

## Home Assistant sample config

![](https://user-images.githubusercontent.com/47263/114361105-753c4200-9b7e-11eb-833c-c26a2b7d0e00.png)

`configuration.yaml`:

```yaml
script:
  air2_sleep:
    icon: mdi:laptop
    sequence:
      - service: mqtt.publish
        data:
          topic: "mac2mqtt/bessarabov-osx/command/sleep"
          payload: "sleep"

  air2_shutdown:
    icon: mdi:laptop
    sequence:
      - service: mqtt.publish
        data:
          topic: "mac2mqtt/bessarabov-osx/command/shutdown"
          payload: "shutdown"

  air2_displaysleep:
    icon: mdi:laptop
    sequence:
      - service: mqtt.publish
        data:
          topic: "mac2mqtt/bessarabov-osx/command/displaysleep"
          payload: "displaysleep"

sensor:
  - platform: mqtt
    name: air2_alive
    icon: mdi:laptop
    state_topic: "mac2mqtt/bessarabov-osx/status/alive"

  - platform: mqtt
    name: "air2_battery"
    icon: mdi:battery-high
    unit_of_measurement: "%"
    state_topic: "mac2mqtt/bessarabov-osx/status/battery"

switch:
  - platform: mqtt
    name: air2_mute
    icon: mdi:volume-mute
    state_topic: "mac2mqtt/bessarabov-osx/status/mute"
    command_topic: "mac2mqtt/bessarabov-osx/command/mute"
    payload_on: "true"
    payload_off: "false"

number:
  - platform: mqtt
    name: air2_volume
    icon: mdi:volume-medium
    state_topic: "mac2mqtt/bessarabov-osx/status/volume"
    command_topic: "mac2mqtt/bessarabov-osx/command/volume"
```

`ui-lovelace.yaml`:

```yaml
title: Home
views:
  - path: default_view
    title: Home
    cards:
      - type: entities
        entities:
          - sensor.air2_alive
          - sensor.air2_battery
          - type: 'custom:slider-entity-row'
            entity: number.air2_volume
            min: 0
            max: 100
          - switch.air2_mute
          - type: button
            name: air2
            entity: script.air2_sleep
            action_name: sleep
            tap_action:
              action: call-service
              service: script.air2_sleep
          - type: button
            name: air2
            entity: script.air2_shutdown
            action_name: shutdown
            tap_action:
              action: call-service
              service: script.air2_shutdown
          - type: button
            name: air2
            entity: script.air2_displaysleep
            action_name: displaysleep
            tap_action:
              action: call-service
              service: script.air2_displaysleep

      - type: history-graph
        hours_to_show: 48
        refresh_interval: 0
        entities:
          - sensor.air2_battery
```

## MQTT topics structure

Program is working with several MQTT topics. All topix are prefixed with `mac2mqtt` + `/` + `COMPUTER_NAME`.
For example, topic with current volume on my machine is `mac2mqtt/bessarabov-osx/status/volume`, in this
case the PREFIX if `mac2mqtt/bessarabov-osx`.

`mac2mqtt` send info to the topics `mac2mqtt/COMPUTER_NAME/status/#` and listen for commands in topics
`mac2mqtt/COMPUTER_NAME/command/#`.

### PREFIX + `/status/alive`

There can be `true` of `false` in this topic. If `mac2mqtt` is connected to MQTT server there is `true`.
If `mac2mqtt` is disconnected from MQTT there is `false`. This is the standard MQTT thing called Last Will and Testament.

### PREFIX + `/status/volume`

The value is the numbers from 0 (inclusive) to 100 (inclusive). The current volume of computer.

The value of this topic is updated every 2 seconds.

### PREFIX + `/status/mute`

There can be `true` of `false` in this topic. `true` means that the computer volume is muted (no sound),
`false` means that it is not multed.

### PREFIX + `/status/battery`

The value is the nuber up to 100. The charge percent of the battery.

The value of this topic is updated every 60 seconds.

### PREFIX + `/command/volume`

You can send integer numberf from 0 (inclusive) to 100 (inclusive) to this topic. It will set the volume on the computer.

### PREFIX + `/command/mute`

You can send `true` of `false` to this topic. When you send `true` the computer is muted. When you send `false` the computer
is unmuted.

### PREFIX + `/command/sleep`

You can send string `sleep` to this topic. It will put computer to sleep mode. Sending some other value will do nothing.

### PREFIX + `/command/shutdown`

You can send string `shutdown` to this topic. It will try to shutdown the computer. The way it is done depends on
the user who run the program. If the program is run by `root` the computer will shutdown, but if it is run by ordinary user
the computer will not shut down if there is other user who logged in.

Sending some other value but `shutdown` will do nothing.

### PREFIX + `/command/displaysleep`

You can send string `displaysleep` to this topic. It will turn off display. Sending some other value will do nothing.
