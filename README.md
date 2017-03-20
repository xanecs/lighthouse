# LightHouse

Control lights and other devices connected to multiple Arduinos with [Redis](https://redis.io) Pub/Sub messaging.

![Build Status](https://travis-ci.org/xanecs/lighthouse.svg?branch=master)

## Installing / Getting started
__Prerequisites__
- A running Redis server
- One or more Arduinos flashed with the _firmata_ sketch.

__Installation__
1. Download a binary for your system from the [releases](https://github.com/xanecs/lighthouse/releases/latest) page.
2. Create a `config.toml` and place it next to the binary (See [configuration](#configuration)).
3. Start the binary

```bash
cd lighthouse
wget https://.../lighthouse
nano config.toml
./lighthouse
```

 Lighthouse should now be connecting to the specified Redis server and await commands on the configured channel.

## Configuration
LightHouse is configured with the `config.toml` file. The service expects this file to be in the _current working directory_.

A minimal `config.toml` might look like this:
```toml
[redis]
host = "127.0.0.1:6379"
topic = "lighthouse"

[boards]
  [boards.first]
  serial = "/dev/ttyUSB0"
    [boards.first.dev]
      [boards.first.dev.led]
      mode = "direct"
      inverted = false
      pinse = ["13"]
```

This tells LightHouse to connect to the redis server running on `127.0.0.1:6379` and subscribe to the channel `lighthouse`.
It configures one Arduino connected to `/dev/ttyUSB0` and exposes the internal LED on pin 13.

#### [redis]
The `[redis]` section is pretty self-explainatory. Specify IP and Port in the `host` field and the Pub/Sub channel in the `topic` field.

#### [boards]
The `[boards]` section is the list of Arduinos to connect to. 

#### [boards.name]
Each Arduino has a `serial` field, specifying the serial port and a `dev` array.

#### [boards.name.dev.name]
Give each device a __unique name across all boards!__ You will later address a device using this name.

__mode__  
what kind of device (for example `"direct"`, `"pwm"` or `"servo"`, see [modes](#modes)).

__inverted__  
some modes have an inverted mode. If you have an LED connected to `+5V` and Arduino pin `2`, you would need to set `inverted` to `true`.

__pins__  
array of pins to use for this mode. some modes require multiple pins (e.g the `rgb` mode requires a red, green and a blue pin.)  
Some modes will also perform actions on all pins if you specify more than one.

## Commands

Commands are received through Redis Pub/Sub. To send a command you publish a JSON string to the configured channel:

```redis
$ redis-cli
redis> PUBLISH lighthouse '{"device": "internal", "action": "on", "params": {}}'
redis> PUBLISH lighthouse '{"device": "internal", "action": "write", "params": {"power": false}}'
```

Each command has the following format
```json
{
  "device": "rgblight",
  "action": "color",
  "params": {
    "red": 255,
    "green": 96,
    "blue": 0
  }
}
```

`device`  
the unique device identifier specified in the `config.toml`

`action`  
which action to perform on the device. supported actiosn depend on the driver

`params`  
parameters for the action

## Modes

Currently, LightHouse supports the following modes:

| Name   | Description              | Inverted          | Pins        |
| ------ | ------------------------ | ----------------- | ----------- |
| direct | switch a pin high or low | `value = !value`  | same on all |
| pwm    | apply pwm to a pin       | `value = 1-value` | same on all |
| rgb    | control an rgb-led       | like pwm          | r, g, b     |
| servo  | control a servo          | does not apply    | same on all |

### Direct

| Action   | Params             |
| -------- | ------------------ |
| `on`     | `{}`               |
| `off`    | `{}`               |
| `write`  | `{"power": bool}`  |

### Pwm

| Action      | Params                                          |
| ------------| ----------------------------------------------- |
| `on`        | `{}`                                            |
| `off`       | `{}`                                            |
| `power`     | `{"power": bool}`                               |
| `brightness`| `{"brightness": float}` (0 to 1)                |
| `write`     | `{"power": bool, "brightness": float}` (0 to 1) |

### Rgb

| Action | Params                                                   |
| -------| -------------------------------------------------------- |
| `on`   | `{}`                                                     |
| `off`  | `{}`                                                     |
| `power`| `{"power": bool}`                                        |
| `color`| `{"red": int, "green": int, "blue": int}` (0 to 255)     |
| `write`| `{"power": bool, "red": int, "green": int, "blue": int}` |

### Servo

| Action | Params                         |
| -------| -------------------------------|
| `on`   | `{}`                           |
| `off`  | `{}`                           |
| `power`| `{"power": bool}`              |
| `angle`| `{"angle": int}` (0 to 180)    |
| `write`| `{"power": bool, "angle": int}`|

## Licensing
The code in this project is licensed under the 2-Clause BSD license (BSD-2-Clause). See `LICENSE` for further details.
