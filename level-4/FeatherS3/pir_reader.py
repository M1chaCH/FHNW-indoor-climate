import board
import digitalio
import time
import proto
import util

pir = None

def init_sensor(): 
    global pir
    pir = digitalio.DigitalInOut(board.D5)
    pir.direction = digitalio.Direction.INPUT

def get_data():
    retry_count = 0
    while True:
        try:
            data = pir.value
            return [proto.Measurement(sensor_type="pir", sensor_value_type=proto.SENSOR_VALUE_TYPE_BOOL, read_timestamp=util.get_timestamp(), sensor_value_name="motion",flag_value=data)]
        except Exception as e:
            retry_count = retry_count + 1
            exception_name = type(e).__name__
            print(f"Failed to load data from pir sensor: {exception_name}: {e}")
            time.sleep(0.5)

            if retry_count > 5:
                raise Exception("Failed to load data from pir sensor") from e
            