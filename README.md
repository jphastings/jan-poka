# Jan Poka

A device for bringing your distant friends a little nearer.

## Set up

1. Set the environment variables in `/lib/systemd/system/jan-poka.service` to your current location, and the direction the device is facing.
2. Reload the system daemon with `sudo systemctl daemon-reload`.
3. Restart the Jan-Poka service with `sudo systemctl restart jan-poka`.
4. Start sending messages to `jan-poka.local`! Here's an example which points at the International Space Station, updating every 15 seconds (that requires an internet connection):
    ```bash
    curl --request GET \
      --url http://jan-poka.local:2678/focus \
      --header 'content-type: application/json' \
      --data '{"poll": 15,"target": [{ "type": "iss" }]}'
   ```

## Options

| $ENV_VAR / buildtag | Requirements | Functionality |
| --- | --- | --- |
| $JP_FACING | - | The direction the device is facing, in degrees clockwise from North |
| $JP_HOMELATITUDE | - | The latitude of where the device is. |
| $JP_HOMELONGITUDE | - | The longitude of where the device is. |
| $JP_HOMEALTITUDE | - | The altitude of where the device is. |
| $JP_USETOWER | A physical tower device | Physically points to the specified locations. | 
| $JP_USEAUDIO | libasound2 | On mac, no dependencies. On linux, [libasound2-dev](https://packages.debian.org/sid/libasound2-dev) library | Allows audio playing (used by text-to-speech) |
| libnova | The [libnova](http://libnova.sourceforge.net/) library | Allows celestial body location |
| rpi | A Raspberry Pi target device | Allows control of stepper motors |
