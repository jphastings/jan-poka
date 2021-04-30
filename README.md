# Jan Poka

A device for bringing your distant friends a little nearer.

## Set up


## Options

| $ENV_VAR / buildtag | Requirements | Functionality |
| --- | --- | --- |
| $JP_FACING | - | The direction the device is facing, in degrees clockwise from North |
| $JP_HOMELATITUDE | - | The latitude of where the device is. |
| $JP_HOMELONGITUDE | - | The longitude of where the device is. |
| $JP_HOMEALTITUDE | - | The altitude of where the device is. |
| $JP_USEAUDIO | libasound2 | On mac, no dependencies. On linux, [libasound2-dev](https://packages.debian.org/sid/libasound2-dev) library | Allows audio playing (used by text-to-speech) |
| $JP_MQTTBROKER | - | The MQTT broker to publish to (host and port). |
| $JP_MQTTUSERNAME | - | The MQTT username to use. |
| $JP_MQTTPASSWORD | - | The MQTT password to use. |
| $JP_MQTTTOPIC | - | The MQTT topic to the target lat/long/alt to. |
| libnova | The [libnova](http://libnova.sourceforge.net/) library | Allows celestial body location |
| rpi | A Raspberry Pi target device | Allows control of stepper motors |
