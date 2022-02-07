from django.http import HttpResponse
from . import settings
from . import myutil
import os
import json
import traceback

def topic_subscribe(req):
    if not req.method == 'POST':
        response = HttpResponse(status=405)
        response.__setitem__('Access-Control-Allow-Origin', '*')
        return response
    try:
        data = json.loads(req.POST['data'])
        print(data)
        response = "sucessfully subscribe"
        if os.path.exists(settings.OUT_PATH):
            response = "already subscribed"
        else:
            myutil.subscribe(data)
        response = HttpResponse(response, content_type="application/json")
        response.__setitem__('Access-Control-Allow-Origin', '*')
        return response
    except:
        traceback.print_exc()
        response = HttpResponse('Bad Request', status=400)
        response.__setitem__('Access-Control-Allow-Origin', '*')
        return response

def ping(request):
    if request.method == "GET":
        response = HttpResponse('sucessful pinging', content_type="application/json")
        response.__setitem__('Access-Control-Allow-Origin', '*')
        return response
