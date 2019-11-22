# Jan Poka

A device for bringing your distant friends a little nearer.



## Options

| $ENV_VAR / buildtag | Requirements | Functionality |
| --- | --- | --- |
| $JP_USEAUDIO | libasound2 | On mac, no dependencies. On linux, [libasound2-dev](https://packages.debian.org/sid/libasound2-dev) library | Allows audio playing (used by text-to-speech) |
| libnova | The [libnova](http://libnova.sourceforge.net/) library | Allows celestial body location |
| rpi | A Raspberry Pi target device | Allows control of stepper motors |
