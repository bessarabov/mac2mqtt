# mac2mqtt

`mac2mqtt` is a program that allow viewing and controlling some aspects of computers running macOS via MQTT.

It publish to MQTT:

 * current volume
 * volume mute state

You can send topics to:

 * change volume
 * mute/unmute
 * put computer to sleep
 * turn off display

## Running

To run this program you need to put 2 files in a directory (`/Users/USERNAME/mac2mqtt/`):

    mac2mqtt
    mac2mqtt.yaml

Edit `mac2mqtt.yaml` (the sample file is in this repository), make binary executable (`chmod +x mac2mqtt`) and run `./mac2mqtt`:

    $ ./mac2mqtt
    2021/04/12 10:37:28 Started
    2021/04/12 10:37:29 Connected to MQTT
    2021/04/12 10:37:29 Sending 'true' to topic: mac2mqtt/bessarabov-osx/status/alive

## Running in the background

You need `mac2mqtt.yaml` and `mac2mqtt` to be placed in the directory `/Users/USERNAME/mac2mqtt/`,
then you need to create file `/Library/LaunchDaemons/com.bessarabov.mac2mqtt.plist`:

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
            entity: script.air2_displaysleep
            action_name: displaysleep
            tap_action:
              action: call-service
              service: script.air2_displaysleep
```

## MQTT topics structure

Program is working with several MQTT topics. All topix are prefixed with `mac2mqtt` + `COMPUTER_NAME`.
For examaple, topic with current volume on my machine is `mac2mqtt/bessarabov-osx/status/volume`

`mac2mqtt` send info to the topics `mac2mqtt/COMPUTER_NAME/status/#` and listen for commands in topics
`mac2mqtt/COMPUTER_NAME/command/#`.

### PREFIX + `/status/alive`

There can be `true` of `false` in this topic. If `mac2mqtt` is connected to MQTT server there is `true`.
If `mac2mqtt` is disconnected from MQTT there is `false`. This is the standard MQTT thing called Last Will and Testament.

### PREFIX + `/status/volume`

The value is the numers from 0 (inclusive) to 100 (inclusive). The current volume of computer.

The value of this topic is updated every 2 seconds.

### PREFIX + `/status/mute`

There can be `true` of `false` in this topic. `true` means that the computer volume is muted (no sound),
`false` means that it is not multed.

### PREFIX + `/command/volume`

You can send integer numberf from 0 (inclusive) to 100 (inclusive) to this topic. It will set the volume on the computer.

### PREFIX + `/command/mute`

You can send `true` of `false` to this topic. When you send `true` the computer is muted. When you send `false` the computer
is unmuted.

### PREFIX + `/command/sleep`

You can send string `sleep` to this topic. It will put computer to sleep mode. Sending some other value will do nothing.

### PREFIX + `/command/displaysleep`

You can send string `displaysleep` to this topic. It will turn off display. Sending some other value will do nothing.
