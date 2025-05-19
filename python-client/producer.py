#!/usr/bin/env python

import os
import json
from dotenv import load_dotenv
from faker import Faker
from confluent_kafka import Producer
from confluent_kafka.admin import AdminClient


def verify_kafka_setup(kafka_config, topic, is_local):    
    if topic is None or topic == "":
        print("⚠️ No topic specified")
        return False
    
    # Local mode still needs to verify topic existence, but no credential check
    if is_local:
        return True
    
    try:
        admin_client = AdminClient(kafka_config)
        metadata = admin_client.list_topics(timeout=5)
        
        if topic not in metadata.topics:
            print(f"Topic '{topic}' not found.")
            return False       
        return True
    except Exception as e:
        # Log the specific exception for easier troubleshooting
        print(f"⚠️ Kafka connection error: {e}")
        return False


def delivery_callback(err, msg):
    if err:
        print(f"ERROR: Message failed delivery: {err}")
    else:
        print(f"Message delivered to topic {msg.topic()} partition [{msg.partition()}] at offset {msg.offset()}")


if __name__ == "__main__":
    load_dotenv()

    topic = os.getenv("CC_TOPIC")
    bootstrap_server = os.getenv("CC_BOOTSTRAP_SERVER")

    # Check for local environments including Docker
    is_local = any(local_indicator in bootstrap_server for local_indicator in 
                  ["localhost", "127.0.0.1", "kafka:"])

    kafka_config = {
        "bootstrap.servers": bootstrap_server,
        "client.id": os.getenv("CLIENT_ID"),
    }

    # Configure security based on environment
    if is_local:
        kafka_config["security.protocol"] = "PLAINTEXT"
        print("🧪 Using PLAINTEXT protocol for local broker")
    else:
        kafka_config.update({
            "security.protocol":        "ssl",
            "ssl.ca.location":          "../ca.pem",
            "ssl.key.location":         "../client-key.pem",
            "ssl.certificate.location": "../client-cert.pem",
            "ssl.endpoint.identification.algorithm": "none",
    		"socket.timeout.ms":        "5000",
            "acks": "all",
            "linger.ms": 20,
            "batch.size": 16384,
            "max.in.flight.requests.per.connection": 5,
            "retries": 3,
            "retry.backoff.ms": 1000,
            "compression.type": "snappy",
            "enable.idempotence": True,
        })
        print("☁️ Using SASL_SSL protocol for cloud broker")

    if verify_kafka_setup(kafka_config, topic, is_local) is False:
        print(f"❌ Kafka configuration error - Exiting")
        raise RuntimeError("Failed to verify Kafka setup")
    print(f"✅ Connected to Kafka ({bootstrap_server})")
    
    # Initialize producer
    producer = Producer(kafka_config)
    fake = Faker()

    serializer = None
    print("⚠️ Using basic JSON serialization (no Schema Registry)")
    
    print(f"Producing to topic '{topic}'...")

    # Generate and send messages
    message_count = 5000
    for i in range(message_count):
        key = fake.uuid4()
        value = {
            "TransactionId": key,
            "AccountNumber": fake.iban(),
            "Amount": round(fake.random_number(digits=5) + fake.random.random(), 2),
            "Currency": fake.currency_code(),
            "Timestamp": fake.iso8601(),
            "TransactionType": fake.random_element(["deposit", "withdrawal", "transfer", "payment"]),
            "Status": fake.random_element(["pending", "completed", "failed"]),
        }

        serialized_value = json.dumps(value).encode('utf-8')

        producer.produce(topic, key=key, value=serialized_value, callback=delivery_callback)
        # producer.poll(100)

        print(f"Producing message {i+1}/{message_count}")
        # print(f"{json.dumps(value)}")

    producer.flush()
    print(f"✅ Successfully produced {message_count} messages to topic {topic}")
    