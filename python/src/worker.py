import asyncio

from temporalio.client import Client
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
    client = await Client.connect("localhost:7233")

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
