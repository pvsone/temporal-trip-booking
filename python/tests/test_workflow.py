import uuid

import pytest
from temporalio import activity
from temporalio.client import Client
from temporalio.worker import Worker

from src.activities import book_flight, book_hotel, book_car, notify_user
from src.shared import BookTripInput
from src.workflow import BookWorkflow


@pytest.fixture
def task_queue() -> str:
    return str(uuid.uuid4())


@pytest.mark.asyncio
async def test_book_workflow_success(client: Client, task_queue: str) -> None:
    input: BookTripInput = BookTripInput(
        userId="Test User",
        flightId="Flight 123",
        hotelId="Hotel 456",
        carId="Car 789"
    )
    async with Worker(
            client,
            task_queue=task_queue,
            workflows=[BookWorkflow],
            activities=[book_flight, book_hotel, book_car, notify_user],
    ):
        result: str = await client.execute_workflow(
            BookWorkflow.run,
            input,
            id=str(uuid.uuid4()),
            task_queue=task_queue,
        )
        assert result == "Booked flight: Flight 123 Booked hotel: Hotel 456 Booked car: Car 789"


@activity.defn(name="book_flight")
async def book_flight_mocked(input: BookTripInput) -> str:
    return f"Booked flight: {input.flightId}"


@activity.defn(name="book_hotel")
async def book_hotel_mocked(input: BookTripInput) -> str:
    return f"Booked hotel: {input.hotelId}"


@activity.defn(name="book_car")
async def book_car_mocked(input: BookTripInput) -> str:
    return f"Booked car: {input.carId}"


@activity.defn(name="notify_user")
async def notify_user_mocked(input: BookTripInput) -> str:
    return f"Notified user: {input.userId}"


@pytest.mark.asyncio
async def test_book_workflow_success_mocked_activities(client: Client, task_queue: str) -> None:
    input: BookTripInput = BookTripInput(
        userId="Test User",
        flightId="Flight 123",
        hotelId="Hotel 456",
        carId="Car 789"
    )
    async with Worker(
            client,
            task_queue=task_queue,
            workflows=[BookWorkflow],
            activities=[book_flight_mocked, book_hotel_mocked, book_car_mocked, notify_user_mocked],
    ):
        result: str = await client.execute_workflow(
            BookWorkflow.run,
            input,
            id=str(uuid.uuid4()),
            task_queue=task_queue,
        )
        assert result == "Booked flight: Flight 123 Booked hotel: Hotel 456 Booked car: Car 789"
