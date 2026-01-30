import asyncio
import adafruit_scd30

import os
import board
import util
import proto
import config

scd = None

def init_sensor(): 
    i2c = board.I2C()
    global scd
    scd = adafruit_scd30.SCD30(i2c)

def update_config():
    scd.measurement_interval = config.get_or_set("SCD30_MEASUREMENT_INTERVAL", os.getenv("SCD30_MEASUREMENT_INTERVAL", 5), proto.CONFIG_OPTION_TYPE_INT32)
    scd.altitude = config.get_or_set("SCD30_ALTITUDE", os.getenv("SCD30_ALTITUDE", 417), proto.CONFIG_OPTION_TYPE_INT32)
    scd.ambient_pressure = config.get_or_set("SCD30_AMBIENT_PRESSURE", os.getenv("SCD30_AMBIENT_PRESSURE", 964), proto.CONFIG_OPTION_TYPE_INT32)
    scd.temperature_offset = config.get_or_set("SCD30_TEMPERATURE_OFFSET", os.getenv("SCD30_TEMPERATURE_OFFSET", 0), proto.CONFIG_OPTION_TYPE_INT32)

async def get_data():
    retry_count = 0
    data = None

    while retry_count < 5 and data == None:
        try:
            if scd.data_available:
                data = [create_proto_measurement("co2", scd.CO2), create_proto_measurement("temp", scd.temperature), create_proto_measurement("hum", scd.relative_humidity)]
            else: 
                await asyncio.sleep(0.2)
        except Exception as e:
            retry_count = retry_count + 1
            exception_name = type(e).__name__
            print(f"Failed to load data from scd30 sensor: {exception_name}: {e}")
            await asyncio.sleep(0.5)

    if data == None:
        raise Exception("Failed to load data from scd30 sensor") from e
    else: 
        return data
            
def create_proto_measurement(key, value):
    return proto.Measurement(sensor_type="scd30", sensor_value_type=proto.SENSOR_VALUE_TYPE_DOUBLE, read_timestamp=util.get_timestamp(), sensor_value_name=key,double_value=value,)
