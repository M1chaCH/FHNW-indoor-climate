import adafruit_minimqtt.adafruit_minimqtt as MQTT
import os
import time
import util

broker_url = os.getenv("BROKER_URL")
broker_port = os.getenv("BROKER_PORT")
broker_username = os.getenv("BROKER_USERNAME")
broker_password = os.getenv("BROKER_PASSWORD")

mqtt_client = None

TOPIC_DATA = "sensor/data"

def init_connection(socket_pool, ssl_context):
    global mqtt_client

    if mqtt_client != None:
        mqtt_client.disconnect()

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

    try_connect_broker()

    
def publish_data(body):
    device_id = util.get_device_id()
    retry_count = 0
    while True:
        try:
            mqtt_client.publish(f"{TOPIC_DATA}/{device_id}", body, False, 0)
            return
        except MQTT.MMQTTException as mqttException: 
            retry_count = retry_count + 1
            print(f"Failed to send data to server: MMQTTException: {mqttException}")
            print("reconnecting...")
            try_connect_broker(True)

        except Exception as e:
            retry_count = retry_count + 1
            exception_name = type(e).__name__
            print(f"Failed to send data to server: {exception_name}: {e}")
            time.sleep(1)

        if retry_count > 5:
            raise Exception("Failed to send data to server") from e
            

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
    print(f"New message on topic {topic}: {message}")


def try_connect_broker(reconnect = False):
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
            time.sleep(1)

            if retry_count > 5:
                raise Exception(f"Failed to connect to mqtt broker at {broker_url}:{broker_port}") from e
