import adafruit_scd30

import os
import board
import util
import time
import proto

scd = None

def init_sensor(): 
    i2c = board.I2C()
    global scd
    scd = adafruit_scd30.SCD30(i2c)
    scd.measurement_interval = os.getenv("SCD30_MEASUREMENT_INTERVAL", 5)
    scd.altitude = os.getenv("SCD30_ALTITUDE", 417)
    scd.ambient_pressure = os.getenv("SCD30_AMBIENT_PRESSURE", 964)
    scd.temperature_offset = os.getenv("SCD30_TEMPERATURE_OFFSET", 0)


def get_data():
    retry_count = 0
    while True:
        try:
            if scd.data_available:
                return [create_proto_measurement("co2", scd.CO2), create_proto_measurement("temp", scd.temperature), create_proto_measurement("hum", scd.relative_humidity)]
            time.sleep(0.2)
        except Exception as e:
            retry_count = retry_count + 1
            exception_name = type(e).__name__
            print(f"Failed to load data from scd30 sensor: {exception_name}: {e}")
            time.sleep(0.5)

            if retry_count > 5:
                raise Exception("Failed to load data from scd30 sensor") from e
            
def create_proto_measurement(key, value):
    return proto.Measurement(sensor_type="scd30", sensor_value_type=proto.SENSOR_VALUE_TYPE_DOUBLE, read_timestamp=util.get_timestamp(), sensor_value_name=key,double_value=value,)