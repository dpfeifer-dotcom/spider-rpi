import requests
from threading import Thread

ENDPOINT = "http://spider-nginx-proxy:9000/"


def switch_light(endpoint: str):
    def _send_request():
        try:
            response = requests.get(f"{ENDPOINT}light/light?type={endpoint}")
            print(f"Request sent, status code: {response.status_code}")
        except Exception as e:
            print(f"Failed to send request: {e}")
    
    Thread(target=_send_request, daemon=True).start()

def switch_bulb():
    def _send_request():
        try:
            response = requests.get(f"{ENDPOINT}light/bulb")
            print(f"Request sent, status code: {response.status_code}")
        except Exception as e:
            print(f"Failed to send request: {e}")
    
    Thread(target=_send_request, daemon=True).start()
