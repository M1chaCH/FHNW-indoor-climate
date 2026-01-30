import asyncio
import adafruit_minimqtt.adafruit_minimqtt as MQTT
import os
import util
import config
import proto

broker_url = os.getenv("BROKER_URL")
broker_port = os.getenv("BROKER_PORT")
broker_username = os.getenv("BROKER_USERNAME")
broker_password = os.getenv("BROKER_PASSWORD")

mqtt_client = None
config_change_restart = False

TOPIC_DATA = "sensor/data"
TOPIC_CONFIG = "device/config"
TOPIC_CONFIG_UPDATE = "device/config/update"

async def init_connection(socket_pool, ssl_context):
    global mqtt_client
    global config_change_restart

    if mqtt_client != None:
        mqtt_client.disconnect()
    config_change_restart = False

    mqtt_client = MQTT.MQTT(
        broker=broker_url,
        port=broker_port,
        username=broker_username,
        password=broker_password,
        socket_pool=socket_pool,
        ssl_context=ssl_context,
        is_ssl=False,
        use_binary_mode=True,
        connect_retries=10,
    )

    mqtt_client.on_connect = on_connect
    mqtt_client.on_disconnect = on_disconnect
    mqtt_client.on_subscribe = on_subscribe
    mqtt_client.on_unsubscribe = on_unsubscribe
    mqtt_client.on_publish = on_publish
    mqtt_client.on_message = on_message

    await try_connect_broker()

def subscribe_config_change():
    mqtt_client.subscribe(f"{TOPIC_CONFIG_UPDATE}/{util.get_device_id()}", 1)

def run_mqtt_listen_loop():
    mqtt_client.loop()

async def publish_data(body):
    device_id = util.get_device_id()
    await try_send(body, f"{TOPIC_DATA}/{device_id}", "sensor data")
        

async def publish_config(body):
    device_id = util.get_device_id()
    await try_send(body, f"{TOPIC_CONFIG}/{device_id}", "device config")
            

async def try_send(body, topic, data_name):
    retry_count = 0
    while True:
        try:
            mqtt_client.publish(topic, body, False, 0)
            return
        except MQTT.MMQTTException as mqttException: 
            retry_count = retry_count + 1
            print(f"Failed to send {data_name} to server: MMQTTException: {mqttException}")
            print("reconnecting...")
            await try_connect_broker(True)

        except Exception as e:
            retry_count = retry_count + 1
            exception_name = type(e).__name__
            print(f"Failed to send {data_name} to server: {exception_name}: {e}")
            await asyncio.sleep(1)

        if retry_count > 5:
            raise Exception("Failed to send {data_name} to server") from e


def on_connect(mqtt_client, userdata, flags, rc):
    print("Connected to MQTT Broker!")
    print(f"Flags: {flags}\n RC: {rc}")


def on_disconnect(mqtt_client, userdata, rc):
    print("Disconnected from MQTT Broker!")


def on_subscribe(mqtt_client, userdata, topic, granted_qos):
    print(f"Subscribed to {topic} with QOS level {granted_qos}")


def on_unsubscribe(mqtt_client, userdata, topic, pid):
    print(f"Unsubscribed from {topic} with PID {pid}")


def on_publish(mqtt_client, userdata, topic, pid):
    print(f"Published to {topic} with PID {pid}")


def on_message(client, topic, message):
    print(f"New message on topic {topic}")

    if topic == f"{TOPIC_CONFIG_UPDATE}/{util.get_device_id()}":
        handleConfigUpdateMessage(message)


async def try_connect_broker(reconnect = False):
    retry_count = 0
    while True:
        try:
            if reconnect: mqtt_client.reconnect()
            else: mqtt_client.connect()
            print(f"successfully connected to broker at {broker_url}:{broker_port}")
            return
        except Exception as e:
            retry_count = retry_count + 1
            exception_name = type(e).__name__
            print(f"Failed to connect to mqtt broker at {broker_url}:{broker_port}: {exception_name}: {e}")
            await asyncio.sleep(1)

            if retry_count > 5:
                raise Exception(f"Failed to connect to mqtt broker at {broker_url}:{broker_port}") from e


def handleConfigUpdateMessage(message):
    global config_change_restart
    configs = proto.DeviceConfigOptions.decode(message)
    config_change_restart = config.set_from_proto(configs)

def has_config_changed():
    return config_change_restart