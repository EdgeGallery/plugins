import time
import logging
import paho.mqtt.client as mqtt
import traceback
import json

import sys
from os import getenv, environ

logging.basicConfig(filename='recv_client.log', level=logging.INFO, filemode='w')


def connect(client_name, broker_adress, port=1883):
    mqttc = mqtt.Client(client_name, True)
    if 'MQTT_EMQX_USER' in environ:
      mqttc.username_pw_set(getenv('MQTT_EMQX_USER'), getenv('MQTT_EMQX_PWD'))
    mqttc.on_connect = on_connect  # attach function to callback
    mqttc.on_message = on_message  # attach function to callback
    mqttc.connect(broker_adress, port)  # connect to broker
    mqttc.loop_start()  # start loop to process callback
    time.sleep(4)  # wait for connection setup to complete

    return mqttc


def on_connect(client, userdata, flags, rc):
    """
    The callback for when the client receives a CONNACK response from the server.
    """
    if rc == 0:
        client.subscribe(topic_subscribe)
        print("connected OK with returned code=", rc)
    else:
        print("Bad connection with returned code=", rc)

def on_message(client, userdata, msg):
    """
    The callback for when a PUBLISH message is received from the server
    """

    message = json.loads(msg.payload)[0]
    print("---------------------------------------------------------------")
    message['ts'] = json.dumps(message['ts'])
    message['data'] = json.dumps(message['data'])
    print(message)

    if message:
        try:
            msg = json.dumps(message)
            client.publish(topic_publish, msg, 0)

	    # Finished message
            print("Finished sending message " )

        except Exception as e:
            print(e)
            traceback.print_tb(e.__traceback__)
            sys.exit()


client_name = "dpsim_adapter"
topic_subscribe = "/fledge"
topic_publish = "/dpsim-fledge"

# Local Fledge platform broker
broker_address_emqx = getenv('MQTT_EMQX_BROKER', 'localhost')
port_emqx = int(getenv('MQTT_EMQX_PORT','1883'))

mqttc = connect(client_name, broker_address_emqx, port_emqx)

print("Press CTRL+C to stop client...")

mqttc.publish("/debug", "SE started")

try:
    while 1:
        time.sleep(1)
        # ensure debug ouput gets flushed
        sys.stdout.flush()

except KeyboardInterrupt:
    print('Exiting...')
    mqttc.loop_stop()
    mqttc.disconnect()
    sys.exit(0)
