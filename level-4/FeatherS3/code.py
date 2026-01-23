import adafruit_connection_manager

import os
import time
import wifi
import mqtt
import proto
import util

scd30_connected = False
try:
    import scd30_reader
    scd30_connected = os.getenv("MODULE_SCD30") == 1
    print("SCD30 module successfully imported", scd30_connected)
except Exception as e:
    print("SCD30 module not found")

pir_connected = False
try:
    import pir_reader
    pir_connected = os.getenv("MODULE_PIR") == 1
    print("PIR module successfully imported", pir_connected)
except Exception as e:
    print("SCD30 module not found")

pm25_connected = False
try:
    import pm25_reader
    pm25_connected = os.getenv("MODULE_PM25") == 1
    print("pm25 module successfully imported", pm25_connected)
except Exception as e:
    print("pm25 module not found")

# load environment variables
ssid = os.getenv("WIFI_SSID")
password = os.getenv("WIFI_PASSWORD")
device_name = os.getenv("DEVICE_NAME")

# Initialize Wifi, Socket Pool, Request Session
pool = adafruit_connection_manager.get_radio_socketpool(wifi.radio)
ssl_context = adafruit_connection_manager.get_radio_ssl_context(wifi.radio)

while True:
    try:
        if scd30_connected:
            scd30_reader.init_sensor()

        if pir_connected:
            pir_reader.init_sensor()

        if pm25_connected:
            pm25_reader.init_sensor()

        print(f"\nConnecting to {ssid}...")
        wifi.radio.connect(
            ssid, password
        )  # automatically tries to reconnect if connection was established
        print("Connection successfully established")

        mqtt.init_connection(pool, ssl_context)

        while True:
            measurements = []

            if scd30_connected:
                measurements.extend(scd30_reader.get_data())

            if pir_connected:
                measurements.extend(pir_reader.get_data())

            if pm25_connected:
                pm25_reader.get_data()

            ip = str(wifi.radio.ipv4_address)
            data = proto.SensorData(device_id=util.get_device_id(), ip=ip, device_name=device_name, measurements=measurements)
            encoded_data = data.encode()
            mqtt.publish_data(encoded_data)

    except Exception as e:
        exception_name = type(e).__name__
        print(
            f"Some unhandled exception occurred (restarting in 5 Seconds): {exception_name}: {e}"
        )
        time.sleep(5)

print("stopped")
