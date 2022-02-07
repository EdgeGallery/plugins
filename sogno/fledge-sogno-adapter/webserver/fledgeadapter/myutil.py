import os

def subscribe(data):
    cmd = 'nohup python3 /home/fledge-plugin.py ' + data['topic'] + ' ' + data['brokerHost'] + ' ' + data['brokerPort']
    print(cmd)
    os.system(cmd)
    #return outname

