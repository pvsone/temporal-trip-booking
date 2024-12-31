import asyncio
import logging

from temporalio import activity
from temporalio.exceptions import ApplicationError

from dataclasses import dataclass

@dataclass
class BookTripInput:
    book_user_id: str
    book_car_id: str
    book_hotel_id: str
    book_flight_id: str

logging.basicConfig(level=logging.INFO)


@activity.defn
async def book_car(input: BookTripInput) -> str:
    activity.logger.info(f"Booking car {input.book_car_id}")

    await asyncio.sleep(1)

    if "flaky" in input.book_car_id.lower():
        # a transient error, which will be retried
        if activity.info().attempt < 6:
            raise ApplicationError("Car booking service is currently unavailable")

    return f"Booked car: {input.book_car_id}"


@activity.defn
async def book_hotel(input: BookTripInput) -> str:
    activity.logger.info(f"Booking hotel: {input.book_hotel_id}")

    await asyncio.sleep(1)

    if "buggy" in input.book_hotel_id.lower():
        # a simulated bug
        error = True
        if error:
            raise Exception("Error due to bug in code")

    return f"Booked hotel: {input.book_hotel_id}"


@activity.defn
async def book_flight(input: BookTripInput) -> str:
    activity.logger.info(f"Booking flight: {input.book_flight_id}")

    await asyncio.sleep(1)

    if "invalid" in input.book_flight_id.lower():
        # a business error, which cannot be retried
        raise ApplicationError(f"Flight {input.book_flight_id} is invalid", type="InvalidFlight", non_retryable=True)

    return f"Booked flight: {input.book_flight_id}"

@activity.defn
async def notify_user(input: BookTripInput) -> str:
    activity.logger.info(f"Notifying user: {input.book_user_id}")

    await asyncio.sleep(1)

    return f"Notified user: {input.book_user_id}"


@activity.defn
async def undo_book_car(input: BookTripInput) -> str:
    activity.logger.info(f"Undo car booking: {input.book_car_id}")

    await asyncio.sleep(1)

    return f"Unbooked car: {input.book_car_id}"


@activity.defn
async def undo_book_hotel(input: BookTripInput) -> str:
    activity.logger.info(f"Undo hotel booking: {input.book_hotel_id}")

    await asyncio.sleep(1)

    return f"Unbooked hotel: {input.book_hotel_id}"


@activity.defn
async def undo_book_flight(input: BookTripInput) -> str:
    activity.logger.info(f"Undo flight booking: {input.book_flight_id}")

    await asyncio.sleep(1)

    return f"Unbooked flight: {input.book_flight_id}"
