import json
import random
import time
from django.http import HttpResponseBadRequest
from rest_framework.response import Response
from rest_framework import status
from rest_framework.decorators import api_view
import requests

from async_service.hamming_code import calc_hamming_code, detect_error, correct_hamming_code

from concurrent import futures

executor = futures.ThreadPoolExecutor(max_workers=1)

SECRET_TOKEN = "rust"
CALLBACK_URL = "http://localhost:8080/api/encryptDecryptRequest/update_calculated"
ERROR_PROBABITY = 0.3


def calc(data: str, encode: bool):
    if random.random() < ERROR_PROBABITY:
        return False
    
    if encode:
        return calc_hamming_code(data)

    error = detect_error(data)
    return correct_hamming_code(data, error)

def logic(req_id, to_be_calcuated):
    time.sleep(5)
    print(to_be_calcuated)
    return {
        "req_id": req_id,
        "calculated": [(data['id'], calc(data['data'], data['encode'])) for data in to_be_calcuated]
    }

def logic_callback_send(task): 
    try:
      result = task.result()
      print(result)
    except futures._base.CancelledError:
      return
    
    answer = [
        {
            "id": res[0],
            "success": res[1] != False,
            "result": (res[1] if res[1] != False else None)
        }
        for res in result["calculated"]
    ]

    answer = {
        "req_id": result["req_id"],
        "calculated": answer,
        "token": SECRET_TOKEN,
    }

    r = requests.put(CALLBACK_URL, json=answer, timeout=3)
    print("go resp:", r.text, r.status_code)

'''
{
    "req_id": 512,
    "calc": [
        {
            "id": 5,
            "encode": true,
            "data": 1515
        },
        {
            "id": 6,
            "encode": false,
            "data": 1515
        }
    ]
}
'''

@api_view(['POST'])
def calculate_view(request, format=None):
    if request.method == 'POST':
        try:
            input = json.loads(request.body)
        except json.JSONDecodeError:
            return HttpResponseBadRequest('Invalid Json')      

        task = executor.submit(logic, input['req_id'], input['calc'])
        task.add_done_callback(logic_callback_send)        

        return Response(status=status.HTTP_200_OK)
    else: 
        return HttpResponseBadRequest('Unsupported Method')
