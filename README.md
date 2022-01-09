# Jan Poka

A device for bringing your distant friends a little nearer.

## Set up

Use docker to run anywhere:

```bash
docker run --env-file .env ghcr.io/jphastings/jan-poka:latest
```

(Assuming that you have appropriate environment variables in a file called `.env`)

## Options

While executing:

| $ENV_VAR            | Functionality                                                                                                                       |
|---------------------|-------------------------------------------------------------------------------------------------------------------------------------|
| $JP_FACING          | The direction the device is facing, in degrees clockwise from North                                                                 |
| $JP_HOMELATITUDE    | The latitude of where the device is.                                                                                                |
| $JP_HOMELONGITUDE   | The longitude of where the device is.                                                                                               |
| $JP_HOMEALTITUDE    | The altitude of where the device is.                                                                                                |
| $JP_MQTTBROKER      | The MQTT broker to publish to (host and port).                                                                                      |
| $JP_MQTTUSERNAME    | The MQTT username to use.                                                                                                           |
| $JP_MQTTPASSWORD    | The MQTT password to use.                                                                                                           |
| $JP_MQTTTOPIC       | The MQTT topic to the target lat/long/alt to.                                                                                       |
| $JP_USEAUDIO        | Allows audio playing — used by text-to-speech. (Requires [libasound2-dev](https://packages.debian.org/sid/libasound2-dev) on Linux) |

While building:

| build tag | Functionality                                                                        |
|-----------|--------------------------------------------------------------------------------------|
| libnova   | Allows celestial body location (Requires [libnova](http://libnova.sourceforge.net/)) |
