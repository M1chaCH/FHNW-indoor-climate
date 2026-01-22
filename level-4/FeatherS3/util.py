import microcontroller
import time

def get_timestamp():
    t = time.localtime()
    return "{:04d}-{:02d}-{:02d}T{:02d}:{:02d}:{:02d}".format(
        t.tm_year, t.tm_mon, t.tm_mday, t.tm_hour, t.tm_min, t.tm_sec
    )


def get_device_id():
    raw_uid = microcontroller.cpu.uid
    return "".join("{:02x}".format(b) for b in raw_uid)