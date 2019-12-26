#!/bin/bash
# /usr/local/bin/autoAP.sh
device=wlan0

configure_ap () {
    if [ -e /etc/systemd/network/08-CLI.network ]; then
        mv /etc/systemd/network/08-CLI.network /etc/systemd/network/08-CLI.network~
        systemctl restart systemd-networkd
    fi
}

configure_client () {
    if [ -e /etc/systemd/network/08-CLI.network~ ] &&  wpa_cli -i$device status | grep -q "mode=station"; then
        mv /etc/systemd/network/08-CLI.network~ /etc/systemd/network/08-CLI.network
        systemctl restart systemd-networkd
    fi
}

reconfigure_wpa_supplicant () {
    sleep "$1"
    if [ "$(wpa_cli -i $device all_sta)" = ""]; then
            wpa_cli -i $device reconfigure
    fi
}

case "$2" in

    # Configure access point if one is created
    AP-ENABLED)
        configure_ap
        reconfigure_wpa_supplicant 2m &
        ;;

    # Configure as client, if connected to some network
    CONNECTED)
        configure_client
        ;;

    # Reconfigure wpa_supplicant to search for your wifi again,
    # if nobody is connected to the ap
    AP-STA-DISCONNECTED)
        reconfigure_wpa_supplicant 20 &
        ;;
esac
