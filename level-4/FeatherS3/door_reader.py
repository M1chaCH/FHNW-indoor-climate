import asyncio
import board
import digitalio
import time
import proto
import util

door = None

def init_sensor(): 
    global door

    if door is not None:
        door.deinit()

    door = digitalio.DigitalInOut(board.D9)
    door.direction = digitalio.Direction.INPUT

def update_config():
    return

async def get_data():
    retry_count = 0
    data = None

    while retry_count < 5 and data == None:
        try:
            value = door.value
            data = [proto.Measurement(sensor_type="door", sensor_value_type=proto.SENSOR_VALUE_TYPE_BOOL, read_timestamp=util.get_timestamp(), sensor_value_name="door_button",flag_value=value)]
        except Exception as e:
            retry_count = retry_count + 1
            exception_name = type(e).__name__
            print(f"Failed to load data from door sensor: {exception_name}: {e}")
        await asyncio.sleep(0.5)
            
    if data == None:
        raise Exception("Failed to load data from door sensor") from e
    else:
        return data
