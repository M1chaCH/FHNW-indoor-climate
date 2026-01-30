import asyncio
import board
import digitalio
import time
import proto
import util

pir = None

def init_sensor(): 
    global pir

    if pir is not None:
        pir.deinit()

    pir = digitalio.DigitalInOut(board.D5)
    pir.direction = digitalio.Direction.INPUT

def update_config():
    return

async def get_data():
    retry_count = 0
    data = None

    while retry_count < 5 and data == None:
        try:
            value = pir.value
            data = [proto.Measurement(sensor_type="pir", sensor_value_type=proto.SENSOR_VALUE_TYPE_BOOL, read_timestamp=util.get_timestamp(), sensor_value_name="motion",flag_value=value)]
        except Exception as e:
            retry_count = retry_count + 1
            exception_name = type(e).__name__
            print(f"Failed to load data from pir sensor: {exception_name}: {e}")
        await asyncio.sleep(0.5)
            
    if data == None:
        raise Exception("Failed to load data from pir sensor") from e
    else:
        return data
