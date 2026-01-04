import asyncio

from temporalio.client import Client
from temporalio.envconfig import ClientConfig
from temporalio.worker import Worker

from activities import (
    book_flight,
    book_hotel,
    book_car,
    notify_user,
    undo_book_flight,
    undo_book_hotel,
    undo_book_car,
)
from workflow import BookWorkflow


async def main():
    connect_config = ClientConfig.load_client_connect_config()
    client = await Client.connect(**connect_config)
    print(f"✅ Client connected to {client.service_client.config.target_host} in namespace '{client.namespace}'")

    worker = Worker(
        client,
        task_queue="trip-task-queue",
        workflows=[BookWorkflow],
        activities=[
            book_flight,
            book_hotel,
            book_car,
            undo_book_flight,
            undo_book_hotel,
            undo_book_car,
            notify_user,
        ],
    )

    print("Python trip booking worker starting...")
    await worker.run()


if __name__ == "__main__":
    asyncio.run(main())
