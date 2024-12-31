import logging
from datetime import timedelta
import asyncio

from temporalio import workflow
from temporalio.common import RetryPolicy

with workflow.unsafe.imports_passed_through():
    from activities import BookTripInput, book_car, book_flight, book_hotel, notify_user

logging.basicConfig(level=logging.INFO)


@workflow.defn
class BookWorkflow:
    @workflow.run
    async def run(self, input: BookTripInput):
        workflow.logger.info(f"Book workflow started, user_id = {input.book_user_id}")

        activity_args = {
            "start_to_close_timeout": timedelta(seconds=5),
            "retry_policy": RetryPolicy(
                initial_interval=timedelta(seconds=1),
                backoff_coefficient=2,
                maximum_interval=timedelta(seconds=30)
            ),
        }

        # saga compensations
        compensations = []

        try:
            # Book Car
            car = await workflow.execute_activity(book_car, input, **activity_args)
            compensations.append("undo_book_car")

            workflow.logger.info("Sleeping for 1 second...")
            await asyncio.sleep(1)

            # Book Hotel
            hotel = await workflow.execute_activity(book_hotel, input, **activity_args)
            compensations.append("undo_book_hotel")

            workflow.logger.info("Sleeping for 1 second...")
            await asyncio.sleep(1)

            # Book Flight
            flight = await workflow.execute_activity(book_flight, input, **activity_args)
            compensations.append("undo_book_flight")

            # Send Notification
            await workflow.execute_activity(notify_user, input, **activity_args)

            return f"{car} {hotel} {flight}"

        except Exception as ex:
            workflow.logger.error("Failed to book trip", ex)
            for compensation in reversed(compensations):
                await workflow.execute_activity(
                    compensation,
                    input,
                    **activity_args,
                )

            return "Booking cancelled"
