import adafruit_connection_manager
import adafruit_requests
import adafruit_scd30

import board
import json
import microcontroller
import os
import time
import wifi

# load environment variables
ssid = os.getenv("WIFI_SSID")
password = os.getenv("WIFI_PASSWORD")
api_key = os.getenv("API_KEY")
data_url = os.getenv("DATA_URL")

# Initialize Wifi, Socket Pool, Request Session
pool = adafruit_connection_manager.get_radio_socketpool(wifi.radio)
ssl_context = adafruit_connection_manager.get_radio_ssl_context(wifi.radio)
requests = adafruit_requests.Session(pool, ssl_context)

# Initialize Sensor
i2c = board.I2C()
scd = adafruit_scd30.SCD30(i2c)
scd.measurement_interval = os.getenv("SENSOR_MEASUREMENT_INTERVAL", 5)
scd.altitude = os.getenv("SENSOR_ALTITUDE", 417)
scd.ambient_pressure = os.getenv("SENSOR_AMBIENT_PRESSURE", 417)
scd.temperature_offset = os.getenv("SENSOR_TEMPERATURE_OFFSET", 0)

def get_timestamp():
    t = time.localtime()
    return "{:04d}-{:02d}-{:02d}T{:02d}:{:02d}:{:02d}".format(
        t.tm_year, t.tm_mon, t.tm_mday, t.tm_hour, t.tm_min, t.tm_sec
    )


def get_device_id():
    raw_uid = microcontroller.cpu.uid
    return "".join("{:02x}".format(b) for b in raw_uid)


def get_data():
    retry_count = 0
    while True:
        try:
            if scd.data_available:
                return {
                    "device": get_device_id(),
                    "co2": scd.CO2,
                    "temp": scd.temperature,
                    "hum": scd.relative_humidity,
                    "uptime": get_timestamp(),
                }
            time.sleep(0.2)
        except Exception as e:
            retry_count = retry_count + 1
            exception_name = type(e).__name__
            print(f"Failed to load data from sensor: {exception_name}: {e}")
            time.sleep(0.5)

            if retry_count > 5:
                raise Exception("Failed to load data from sensor") from e


def send_data(body):
    retry_count = 0
    while True:
        try:
            with requests.post(data_url, data=body, headers={"X-Api-Key": api_key, "Content-Type": "application/json"}) as response:
                if response.status_code != 200:
                    retry_count = retry_count + 1
                    if retry_count > 10:
                        raise Exception("Server responded with status code after multiple tries: " + str(response.status_code))
                    else:
                        print(f"WARN: bad response: '{response.status_code}' retrying...")
                        time.sleep(1)

                else: return
        except Exception as e:
            retry_count = retry_count + 1
            exception_name = type(e).__name__
            print(f"Failed to send data to server: {exception_name}: {e}")
            time.sleep(1)

            if retry_count > 5:
                raise Exception("Failed to send data to server") from e


while True:
    try:
        print(f"\nConnecting to {ssid}...")
        wifi.radio.connect(
            ssid, password
        )  # automatically tries to reconnect if connection was established
        print("Connection successfully established")

        while True:
            data = get_data()
            json_string = json.dumps(data)
            print(json_string)

            send_data(json_string)

    except Exception as e:
        exception_name = type(e).__name__
        print(
            f"Some unhandled exception occurred (restarting in 5 Seconds): {exception_name}: {e}"
        )
        time.sleep(5)

print("stopped")
