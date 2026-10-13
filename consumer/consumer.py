"""
South Park message consumer.

Listens on the RabbitMQ queue "southpark_messages" and prints every
message the Go API publishes.
"""

import json
import os
import sys
import time

import pika

# Must match the Go publisher EXACTLY (see go-api/adapters/rabbitmq/publisher.go)
RABBITMQ_URL = os.getenv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/")
EXCHANGE_NAME = "southpark"
QUEUE_NAME = "southpark_messages"
BINDING_KEY = "messages.#"


def connect_with_retry(url, attempts=10, wait_seconds=3):
    """RabbitMQ in Docker needs a few seconds to start, so try a few times."""
    for attempt in range(1, attempts + 1):
        try:
            return pika.BlockingConnection(pika.URLParameters(url))
        except pika.exceptions.AMQPConnectionError:
            print(
                f"RabbitMQ: connection attempt {attempt}/{attempts} failed, "
                f"retrying in {wait_seconds}s...",
                flush=True,
            )
            time.sleep(wait_seconds)
    sys.exit("Could not connect to RabbitMQ, giving up.")


def on_message(channel, method, properties, body):
    """Called by pika for every message that arrives in the queue."""
    try:
        message = json.loads(body)
        author = message.get("author", "Someone")
        text = message.get("body", "")
        print(f'{author}: "{text}"   (routing key: {method.routing_key})', flush=True)
    except json.JSONDecodeError:
        print(f"Received a message that is not JSON: {body!r}", flush=True)

    # Tell RabbitMQ we are done with this message, so it can be removed from the queue
    channel.basic_ack(delivery_tag=method.delivery_tag)


def main():
    # Step 1: connect and open a channel
    connection = connect_with_retry(RABBITMQ_URL)
    channel = connection.channel()

    # Step 2: declare the same exchange, queue and binding as the Go publisher.
    # Declaring is idempotent: if they already exist, nothing changes.
    channel.exchange_declare(exchange=EXCHANGE_NAME, exchange_type="topic", durable=True)
    channel.queue_declare(queue=QUEUE_NAME, durable=True)
    channel.queue_bind(queue=QUEUE_NAME, exchange=EXCHANGE_NAME, routing_key=BINDING_KEY)

    # Step 3: register the consumer (one message at a time)
    channel.basic_qos(prefetch_count=1)
    channel.basic_consume(queue=QUEUE_NAME, on_message_callback=on_message)

    print(" [*] Waiting for messages from South Park. To exit press CTRL+C", flush=True)

    # Step 4: process messages forever
    try:
        channel.start_consuming()
    except KeyboardInterrupt:
        channel.stop_consuming()
    connection.close()


if __name__ == "__main__":
    main()