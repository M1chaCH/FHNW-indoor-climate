from adafruit_pm25.i2c import PM25_I2C

import os
import board
import util
import time
import proto

pm25 = None

def init_sensor(): 
    i2c = board.I2C()
    global pm25
    pm25 = PM25_I2C(i2c, None)

def get_data():
    retry_count = 0
    while True:
        try:
            data = pm25.read()
            print(data)
            return data
        except Exception as e:
            retry_count = retry_count + 1
            exception_name = type(e).__name__
            print(f"Failed to load data from pm25 sensor: {exception_name}: {e}")
            time.sleep(0.5)

            if retry_count > 5:
                raise Exception("Failed to load data from pm25 sensor") from e
            
# def create_proto_measurement(key, value):
#     return proto.SensorData.Measurement(sensor_type="scd30", sensor_value_type=proto.SENSOR_VALUE_TYPE_DOUBLE, read_timestamp=util.get_timestamp(), sensor_value_name=key,double_value=value,)