"""
South Park Spammer - a timed script that sends a random South Park
message to the Go API (POST /messages) every few seconds.
"""

import json
import os
import random
import time
import urllib.error
import urllib.request

# Settings (can be changed with environment variables)
GO_API_URL = os.getenv("GO_API_URL", "http://localhost:8080")
INTERVAL_SECONDS = float(os.getenv("SPAM_INTERVAL_SECONDS", "5"))

# A list of tuples: (author, body)
SOUTH_PARK_MESSAGES = [
    ("Cartman", "Respect my authoritah!"),
    ("Cartman", "Screw you guys, I'm going home!"),
    ("Cartman", "I'm not fat, I'm big-boned!"),
    ("Kyle", "You know, I learned something today..."),
    ("Stan", "Oh my God, they killed Kenny!"),
    ("Kyle", "You bastards!"),
    ("Kenny", "Mmmph mmmph mmph mmmph!"),
    ("Butters", "Oh, hamburgers."),
    ("Randy", "I thought this was America!"),
    ("Mr. Mackey", "Drugs are bad, mmkay?"),
    ("Mr. Mackey", "Spamming the message queue is bad, mmkay?"),
    ("Mr. Garrison", "Alright class, today we learn about message queues."),
    ("Towelie", "Don't forget to bring a towel!"),
    ("Mr. Hankey", "Howdy ho!"),
    ("Timmy", "Timmy!"),
]


def send_message(author, body):
    """POST one message to the Go API and return the HTTP status code."""
    data = json.dumps({"author": author, "body": body}).encode("utf-8")
    request = urllib.request.Request(
        f"{GO_API_URL}/messages",
        data=data,
        headers={"Content-Type": "application/json"},
        method="POST",
    )
    with urllib.request.urlopen(request, timeout=5) as response:
        return response.status


def main():
    print(f"Spammer started: 1 message every {INTERVAL_SECONDS}s to {GO_API_URL}", flush=True)

    while True:
        author, body = random.choice(SOUTH_PARK_MESSAGES)
        try:
            status = send_message(author, body)
            print(f'[{status}] sent {author}: "{body}"', flush=True)
        except urllib.error.HTTPError as error:
            print(f"[{error.code}] the API refused the message", flush=True)
        except (urllib.error.URLError, OSError) as error:
            print(f"Go API not reachable yet ({error}), trying again...", flush=True)

        time.sleep(INTERVAL_SECONDS)


if __name__ == "__main__":
    main()