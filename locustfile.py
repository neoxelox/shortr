from random import randint
from typing import Dict, Tuple

from locust import task
from locust.contrib.fasthttp import FastHttpUser

HEADERS: Dict[str, str] = {"Content-Type": "application/json"}
INTERVAL: Tuple[int, int] = (0, 4096)


class LoadTest(FastHttpUser):
    @task(50)
    def get_url(self):
        with self.client.get(
            f"/{randint(*INTERVAL)}", allow_redirects=False, catch_response=True, headers=HEADERS
        ) as response:
            if response.status_code not in [307, 404]:
                response.failure(f"Got wrong response: {response.status_code}")
            else:
                response.success()

    @task(20)
    def shorten_url(self):
        with self.client.post(
            f"/?url={randint(*INTERVAL)}", allow_redirects=False, catch_response=True, headers=HEADERS
        ) as response:
            if response.status_code not in [200, 400]:
                response.failure(f"Got wrong response: {response.status_code}")
            else:
                response.success()

    @task(5)
    def delete_url(self):
        with self.client.delete(
            f"/{randint(*INTERVAL)}", allow_redirects=False, catch_response=True, headers=HEADERS
        ) as response:
            if response.status_code not in [200, 400]:
                response.failure(f"Got wrong response: {response.status_code}")
            else:
                response.success()

    @task(10)
    def modify_url(self):
        with self.client.put(
            f"/{randint(*INTERVAL)}?url={randint(*INTERVAL)}",
            allow_redirects=False,
            catch_response=True,
            headers=HEADERS,
        ) as response:
            if response.status_code not in [200, 400]:
                response.failure(f"Got wrong response: {response.status_code}")
            else:
                response.success()

    @task(15)
    def get_url_stats(self):
        with self.client.get(
            f"/{randint(*INTERVAL)}/stats", allow_redirects=False, catch_response=True, headers=HEADERS
        ) as response:
            if response.status_code not in [200, 404]:
                response.failure(f"Got wrong response: {response.status_code}")
            else:
                response.success()
