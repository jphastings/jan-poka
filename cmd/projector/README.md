# Projector

Designed for use on the Raspberry Pi (specifically the RPi 0w).

## Notes

Stop the TTY from using the framebuffer:

```bash
# Stop TTY1 from overwriting the framebuffer
sudo systemctl stop getty@tty1.service
# Make the cursor invisible
sudo sh -c '/usr/bin/tput civis > /dev/tty1'
```
