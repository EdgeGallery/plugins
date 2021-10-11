# -*- coding: utf-8 -*-

# FLEDGE_BEGIN
# See: http://fledge-iot.readthedocs.io/
# FLEDGE_END

""" MQTT North plugin"""

# import aiohttp
import asyncio
import json
import paho.mqtt.client as mqtt

from fledge.common import logger
from fledge.plugins.north.common.common import *

__author__ = "Libu Jacob Varghese, Khemendra Kumar"
__copyright__ = "Copyright (c) 2021 Huawei Technologies"
__license__ = "Apache 2.0"
__version__ = "${VERSION}"

_LOGGER = logger.setup(__name__)


mqtt_north = None
config = ""

_CONFIG_CATEGORY_NAME = "MQTT"
_CONFIG_CATEGORY_DESCRIPTION = "MQTT North Plugin"

_DEFAULT_CONFIG = {
    'plugin': {
         'description': 'MQTT North Plugin',
         'type': 'string',
         'default': 'mqtt_north',
         'readonly': 'true'
    },
    'brokerHost': {
        'description': 'Hostname or IP address of the broker to connect to',
        'type': 'string',
        'default': 'mqtt-broker',
        'order': '1',
        'displayName': 'MQTT Broker host',
        'mandatory': 'true'
    },
    'brokerPort': {
        'description': 'The network port of the broker to connect to',
        'type': 'integer',
        'default': '1883',
        'order': '2',
        'displayName': 'MQTT Broker Port',
        'mandatory': 'true'
    },
    'keepAliveInterval': {
        'description': 'Maximum period in seconds allowed between communications with the broker. If no other messages are being exchanged, '
                        'this controls the rate at which the client will send ping messages to the broker.',
        'type': 'integer',
        'default': '60',
        'order': '3',
        'displayName': 'Keep Alive Interval'
    },
    'topic': {
        'description': 'The topic to publish the messages',
        'type': 'string',
        'default': 'Room2/conditions',
        'order': '4',
        'displayName': 'Topic To Publish',
        'mandatory': 'true'
    },
    'qos': {
        'description': 'The desired quality of service level for publish',
        'type': 'integer',
        'default': '0',
        'order': '5',
        'displayName': 'QoS Level',
        'minimum': '0',
        'maximum': '2'
    }
}


def plugin_info():
    return {
        'name': 'mqtt',
        'version': '1.9.1',
        'type': 'north',
        'interface': '1.0',
        'config': _DEFAULT_CONFIG
    }


def plugin_init(data):
    _LOGGER.error("MQTT Broker initialize: ")
    global mqtt_north, config
    mqtt_north = MqttNorthClient(data)
    config = data
    return config


async def plugin_send(data, payload, stream_id):
    # stream_id (log?)
    try:
        is_data_sent, new_last_object_id, num_sent = await mqtt_north.send_payloads(payload)
    except asyncio.CancelledError:
        pass
    else:
        return is_data_sent, new_last_object_id, num_sent


def plugin_shutdown(data):
    mqtt_north.stop()


# TODO: North plugin can not be reconfigured? (per callback mechanism)
def plugin_reconfigure():
    pass


# class HttpNorthPlugin(object):
#     """ North HTTP1 Plugin """

#     def __init__(self):
#         self.event_loop = asyncio.get_event_loop()

#     async def send_payloads(self, payloads):
#         is_data_sent = False
#         last_object_id = 0
#         num_sent = 0
#         try:
#             payload_block = list()

#             for p in payloads:
#                 last_object_id = p["id"]
#                 read = dict()
#                 read["asset"] = p['asset_code']
#                 read["readings"] = p['reading']
#                 read["timestamp"] = p['user_ts']
#                 payload_block.append(read)

#             num_sent = await self._send_payloads(payload_block)
#             is_data_sent = True
#         except Exception as ex:
#             _LOGGER.exception("Data could not be sent, %s", str(ex))

#         return is_data_sent, last_object_id, num_sent

#     async def _send_payloads(self, payload_block):
#         """ send a list of block payloads"""

#         num_count = 0
#         try:
#             verify_ssl = False if config["verifySSL"]['value'] == 'false' else True
#             url = config['url']['value']
#             connector = aiohttp.TCPConnector(verify_ssl=verify_ssl)
#             async with aiohttp.ClientSession(connector=connector) as session:
#                 result = await self._send(url, payload_block, session)
#         except:
#             pass
#         else: 
#             num_count += len(payload_block)
#         return num_count

#     async def _send(self, url, payload, session):
#         """ Send the payload, using provided socket session """
#         headers = {'content-type': 'application/json'}
#         async with session.post(url, data=json.dumps(payload), headers=headers) as resp:
#             result = await resp.text()
#             status_code = resp.status
#             if status_code in range(400, 500):
#                 _LOGGER.error("Bad request error code: %d, reason: %s", status_code, resp.reason)
#                 raise Exception
#             if status_code in range(500, 600):
#                 _LOGGER.error("Server error code: %d, reason: %s", status_code, resp.reason)
#                 raise Exception
#             return result


class MqttNorthClient(object):
    """ mqtt publisher class"""

    def __init__(self, config):
        _LOGGER.error("MQTT __init__ ")
        self.event_loop = asyncio.get_event_loop()
        self.mqtt_client = mqtt.Client()
        self.broker_host = config['brokerHost']['value']
        self.broker_port = int(config['brokerPort']['value'])
        self.topic = config['topic']['value']
        self.qos = int(config['qos']['value'])
        self.keep_alive_interval = int(config['keepAliveInterval']['value'])
        self.start()

    def _write(self, line):
        with open('/mqtt.log', 'a') as f:
            f.write(line + '\n')

    def on_connect(self, client, userdata, flags, rc):
        """ The callback for when the client receives a CONNACK response from the server
        """
        if rc != 0:
            print("Unable to connect to MQTT Broker...")
            _LOGGER.error("Unable to connect to MQTT Broker...")
        else:
            client.connected_flag = True
            _LOGGER.error("Connected with MQTT Broker: ", str(self.broker_host))

    def on_disconnect(self, client, userdata, rc):
        self._write("MQTT disconnected")
        self.start()

    def on_publish(client, userdata, mid):
        pass    

    def start(self):
        _LOGGER.error("MQTT start")
        self._write("MQTT start")
        # event callbacks
        self.mqtt_client.on_connect = self.on_connect
        self.mqtt_client.on_publish = self.on_publish
        self.mqtt_client.on_disconnect = self.on_disconnect
        try:
            self.mqtt_client.connect(self.broker_host, self.broker_port, self.keep_alive_interval)
            _LOGGER.info("MQTT connecting..., Broker Host: %s, Port: %s", self.broker_host, self.broker_port)
            self._write("MQTT connecting..., Broker Host: {}, Port: {}".format(self.broker_host, self.broker_port))
        except ConnectionRefusedError as e:
            _LOGGER.exception(e)
            self._write("MQTT connect failed")

    def stop(self):
        self.mqtt_client.disconnect()

    async def send_payloads(self, payloads):
        self._write("MQTT send_payloads")
        is_data_sent = False
        last_object_id = 0
        num_sent = 0
        for p in payloads:
            last_object_id = p["id"]
            read = dict()
            read["asset"] = p['asset_code']
            read["readings"] = p['reading']
            read["timestamp"] = p['user_ts']
            self._write("Payload: {}".format(read))
            if not await self._send(read):
                return is_data_sent, last_object_id, num_sent
            num_sent = num_sent + 1
            is_data_sent = True

        return is_data_sent, last_object_id, num_sent

    async def _send(self, payload) -> bool:
        """ send a list of block payloads"""
        try:
            info = self.mqtt_client.publish(self.topic, json.dumps(payload))
            if info.rc != mqtt.MQTT_ERR_SUCCESS:
                self._write("Send failed: {}".format(info.rc))
                return False
        except Exception as ex:
            _LOGGER.exception("Data could not be sent, %s", str(ex))
            self._write("Send exception: {}".format(str(ex)))
            return False
        self._write("Send success.")
        return True
