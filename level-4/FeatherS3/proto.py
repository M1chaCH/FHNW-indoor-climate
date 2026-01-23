import minipb

SENSOR_VALUE_TYPE_STRING = 0
SENSOR_VALUE_TYPE_INT32 = 1
SENSOR_VALUE_TYPE_DOUBLE = 2
SENSOR_VALUE_TYPE_BOOL = 3

@minipb.process_message_fields
class Measurement(minipb.Message):
    sensor_type = minipb.Field(1, minipb.TYPE_STRING, required=True)
    sensor_value_type = minipb.Field(2, minipb.TYPE_INT, required=True)
    read_timestamp = minipb.Field(3, minipb.TYPE_STRING)
    sensor_value_name = minipb.Field(4, minipb.TYPE_STRING)
    string_value = minipb.Field(10, minipb.TYPE_STRING)
    int32_value = minipb.Field(11, minipb.TYPE_INT32)
    double_value = minipb.Field(12, minipb.TYPE_DOUBLE)
    flag_value = minipb.Field(13, minipb.TYPE_BOOL)

@minipb.process_message_fields
class SensorData(minipb.Message):
    device_id = minipb.Field(1, minipb.TYPE_STRING, required=True)
    ip = minipb.Field(2, minipb.TYPE_STRING)
    device_name = minipb.Field(3, minipb.TYPE_STRING)

    measurements = minipb.Field(4, Measurement, repeated=True)