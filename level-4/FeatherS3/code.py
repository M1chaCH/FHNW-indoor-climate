import adafruit_connection_manager

import os
import time
import wifi
import scd30_reader
import mqtt

# load environment variables
ssid = os.getenv("WIFI_SSID")
password = os.getenv("WIFI_PASSWORD")
device_name = os.getenv("DEVICE_NAME")

# Initialize Wifi, Socket Pool, Request Session
pool = adafruit_connection_manager.get_radio_socketpool(wifi.radio)
ssl_context = adafruit_connection_manager.get_radio_ssl_context(wifi.radio)

scd30_reader.init_sensor()

while True:
    try:
        print(f"\nConnecting to {ssid}...")
        wifi.radio.connect(
            ssid, password
        )  # automatically tries to reconnect if connection was established
        print("Connection successfully established")

        mqtt.init_connection(pool, ssl_context)

        while True:
            data = scd30_reader.get_data(str(wifi.radio.ipv4_address), device_name)
            encoded_data = data.encode()
            mqtt.publish_data(encoded_data)

    except Exception as e:
        exception_name = type(e).__name__
        print(
            f"Some unhandled exception occurred (restarting in 5 Seconds): {exception_name}: {e}"
        )
        time.sleep(5)

print("stopped")
