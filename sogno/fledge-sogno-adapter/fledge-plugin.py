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

    message = json.loads(msg.payload)
    print("---------------------------------------------------------------")
    readings = message['readings']
    readings['ts'] = json.loads("".join(readings["ts"].split("\\")))
    readings['data'] = json.loads(readings["data"])
    readings = [readings]
    print(readings)

    if readings:
        try:
            sogno_msg = json.dumps(readings)
            #mqttc_sogno.connect(broker_address_rabbitmq, port_rabbitmq)
            mqttc_sogno.publish(topic_publish, sogno_msg, 0)

	    # Finished message
            print("Finished sending message " )

        except Exception as e:
            print(e)
            traceback.print_tb(e.__traceback__)
            sys.exit()

if len(sys.argv) != 4:
    print("Specify three argument")
    sys.exit()

client_name = "sogno_fledge_adapter"
topic_subscribe = sys.argv[1]
topic_publish = "/dpsim-powerflow"

# Local SOGNO platform broker
broker_address_rabbitmq = getenv('MQTT_RABBITMQ_BROKER', 'localhost')
port_rabbitmq = int(getenv('MQTT_RABBITMQ_PORT','1883'))

# Local Fledge platform broker
broker_address_emqx = sys.argv[2]
port_emqx = int(sys.argv[3])

mqttc_sogno = mqtt.Client()
if 'MQTT_RABBITMQ_USER' in environ:
      mqttc_sogno.username_pw_set(getenv('MQTT_RABBITMQ_USER'), getenv('MQTT_RABBITMQ_PWD'))
mqttc_sogno.connect(broker_address_rabbitmq, port_rabbitmq)
mqttc_sogno.loop_start()

mqttc_fledge = connect(client_name, broker_address_emqx, port_emqx)

print("Press CTRL+C to stop client...")

mqttc_fledge.publish("/debug", "SE started")

try:
    while 1:
        time.sleep(1)
        # ensure debug ouput gets flushed
        sys.stdout.flush()

except KeyboardInterrupt:
    print('Exiting...')
    mqttc_fledge.loop_stop()
    mqttc_fledge.disconnect()
    mqttc_sogno.loop_stop()
    mqttc_sogno.disconnect()
    sys.exit(0)
