import asyncio
import logging

from temporalio import activity
from temporalio.exceptions import ApplicationError

from shared import BookTripInput

logging.basicConfig(level=logging.INFO)


@activity.defn
async def book_flight(input: BookTripInput) -> str:
    activity.logger.info(f"Booking flight: {input.flightId}")

    await asyncio.sleep(1)

    if "flaky" in input.flightId.lower():
        # a transient error, which will be retried
        if activity.info().attempt < 6:
            raise ApplicationError("Flight booking service is currently unavailable")

    return f"Booked flight: {input.flightId}"


@activity.defn
async def book_hotel(input: BookTripInput) -> str:
    activity.logger.info(f"Booking hotel: {input.hotelId}")

    await asyncio.sleep(1)

    if "buggy" in input.hotelId.lower():
        # a simulated bug
        error = True
        if error:
            raise Exception("Error due to bug in code")

    return f"Booked hotel: {input.hotelId}"


@activity.defn
async def book_car(input: BookTripInput) -> str:
    activity.logger.info(f"Booking car {input.carId}")

    await asyncio.sleep(1)

    if "invalid" in input.carId.lower():
        # a business error, which cannot be retried
        raise ApplicationError(f"Car {input.carId} is invalid", type="InvalidCar", non_retryable=True)

    return f"Booked car: {input.carId}"


@activity.defn
async def notify_user(input: BookTripInput) -> str:
    activity.logger.info(f"Notifying user: {input.userId}")

    await asyncio.sleep(1)

    return f"Notified user: {input.userId}"


@activity.defn
async def undo_book_flight(input: BookTripInput) -> str:
    activity.logger.info(f"Undo flight booking: {input.flightId}")

    await asyncio.sleep(1)

    return f"Unbooked flight: {input.flightId}"


@activity.defn
async def undo_book_hotel(input: BookTripInput) -> str:
    activity.logger.info(f"Undo hotel booking: {input.hotelId}")

    await asyncio.sleep(1)

    return f"Unbooked hotel: {input.hotelId}"


@activity.defn
async def undo_book_car(input: BookTripInput) -> str:
    activity.logger.info(f"Undo car booking: {input.carId}")

    await asyncio.sleep(1)

    return f"Unbooked car: {input.carId}"
