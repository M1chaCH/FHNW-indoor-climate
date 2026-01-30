import adafruit_connection_manager

import os
import wifi
import mqtt
import proto
import util
import asyncio
import config
import traceback


scd30_connected = False
pir_connected = False
pm25_connected = False
door_connected = False
device_name = ""
running = True
scd30_reader_mod = None
pir_reader_mod = None
pm25_reader_mod = None
door_reader_mod = None


async def run_sensor_loop():
    while running:
        measurements = []

        if scd30_connected:
            measurements.extend(await scd30_reader_mod.get_data())

        if pir_connected:
            measurements.extend(await pir_reader_mod.get_data())

        if pm25_connected:
            measurements.extend(await pm25_reader_mod.get_data())

        if door_connected:
            measurements.extend(await door_reader_mod.get_data())

        ip = str(wifi.radio.ipv4_address)
        data = proto.SensorData(device_id=util.get_device_id(), ip=ip, device_name=device_name, measurements=measurements)
        print("encoding and sending data", data)
        encoded_data = data.encode()
        await mqtt.publish_data(encoded_data)


async def run_push_config_loop():
    print("running config loop")
    while running:
        options = config.get_proto_options()
        data = proto.DeviceConfigOptions(device_id=util.get_device_id(), options=options)
        encoded_config = data.encode()
        await mqtt.publish_config(encoded_config)

        await asyncio.sleep(1 * 60) # run every 1 minutes


async def listen_for_config_change():
    global running
    mqtt.subscribe_config_change()
    while running:
        try:
            mqtt.run_mqtt_listen_loop()
            if mqtt.has_config_changed():
                print("config changed, restarting")
                running = False
            else:
                await asyncio.sleep(1)
        except Exception as e:
            traceback.print_exception(e)
            print("error while listening to mqtt messages")
            await asyncio.sleep(5)


async def main():
    global scd30_connected
    global pir_connected
    global pm25_connected
    global door_connected
    global device_name
    global scd30_reader_mod
    global pir_reader_mod
    global pm25_reader_mod
    global door_reader_mod
    global running

    sensorTask = None
    pushConfigTask = None
    configChangeListenerTask = None

    while True:
        try:
            running = True
            try:
                import scd30_reader as scd30_reader_mod
                scd30_connected = config.get_or_set("MODULE_SCD30", os.getenv("MODULE_SCD30", 0) == 1, proto.CONFIG_OPTION_TYPE_BOOL)
                print("SCD30 module successfully imported", scd30_connected)
            except Exception as e:
                print("SCD30 module not found")
                traceback.print_exception(e)
            
            try:
                import pir_reader as pir_reader_mod
                pir_connected = config.get_or_set("MODULE_PIR", os.getenv("MODULE_PIR", 0) == 1, proto.CONFIG_OPTION_TYPE_BOOL)
                print("PIR module successfully imported", pir_connected)
            except Exception as e:
                print("PIR module not found")
                traceback.print_exception(e)

            try:
                import pm25_reader as pm25_reader_mod
                pm25_connected = config.get_or_set("MODULE_PM25", os.getenv("MODULE_PM25", 0) == 1, proto.CONFIG_OPTION_TYPE_BOOL)
                print("pm25 module successfully imported", pm25_connected)
            except Exception as e:
                print("pm25 module not found")
                traceback.print_exception(e)

            try:
                import door_reader as door_reader_mod
                door_connected = config.get_or_set("MODULE_DOOR", os.getenv("MODULE_DOOR", 0) == 1, proto.CONFIG_OPTION_TYPE_BOOL)
                print("door button module successfully imported", door_connected)
            except Exception as e:
                print("door button module not found")
                traceback.print_exception(e)
    
            if scd30_connected is False and pir_connected is False and pm25_connected is False and door_connected is False:
                print("no sensor connected, quitting")
                return

            # load environment variables
            ssid = os.getenv("WIFI_SSID")
            password = os.getenv("WIFI_PASSWORD")
            device_name = config.get_or_set("DEVICE_NAME", os.getenv("DEVICE_NAME", f"ESP32 - {util.get_device_id()}"), proto.CONFIG_OPTION_TYPE_STRING)

            # Initialize Wifi, Socket Pool, Request Session
            pool = adafruit_connection_manager.get_radio_socketpool(wifi.radio)
            ssl_context = adafruit_connection_manager.get_radio_ssl_context(wifi.radio)

            if scd30_connected:
                scd30_reader_mod.init_sensor()
                scd30_reader_mod.update_config()

            if pir_connected:
                pir_reader_mod.init_sensor()
                pir_reader_mod.update_config()

            if pm25_connected:
                pm25_reader_mod.init_sensor()
                pm25_reader_mod.update_config()

            if door_connected:
                door_reader_mod.init_sensor()
                door_reader_mod.update_config()

            print(f"\nConnecting to {ssid}...")
            wifi.radio.connect(
                ssid, password
            )  # automatically tries to reconnect if connection was established
            print("Connection successfully established")

            await mqtt.init_connection(pool, ssl_context)

            sensorTask = asyncio.create_task(run_sensor_loop())
            pushConfigTask = asyncio.create_task(run_push_config_loop())
            configChangeListenerTask = asyncio.create_task(listen_for_config_change())
            await asyncio.gather(sensorTask, pushConfigTask, configChangeListenerTask)

            print("all tasks complete, restarting...")

        except Exception as e:
            traceback.print_exception(e)
            
            if sensorTask is not None:
                sensorTask.cancel()
            
            if pushConfigTask is not None:
                pushConfigTask.cancel()
                
            print("Some unhandled exception occurred (restarting in 5 Seconds)")
            await asyncio.sleep(5)

asyncio.run(main())
print("done")