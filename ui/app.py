import asyncio
import uuid

from quart import Quart, render_template, request, Response
from temporalio.client import WorkflowExecutionStatus

from client import get_client
from data import BookTripInput

app = Quart(__name__)

client = None


@app.before_serving
async def startup():
    global client
    client = await get_client()


@app.route("/")
async def display_form():
    return await render_template("index.html")


@app.route("/book_progress", methods=["POST"])
async def book_progress():
    form = await request.form
    user_id = f'{form.get("name").replace(
        " ", "-").lower()}-{str(uuid.uuid4().int)[:6]}'
    flight = form.get("flight")
    hotel = form.get("hotel")
    car = form.get("car")

    trip_input = BookTripInput(
        userId=user_id,
        flightId=flight,
        hotelId=hotel,
        carId=car,
    )

    # Start workflow without waiting for result
    await client.start_workflow(
        "BookWorkflow",
        trip_input,
        id=user_id,
        task_queue="trip-task-queue",
    )

    return await render_template("book_progress.html", workflow_id=user_id)


@app.route("/progress_stream/<workflow_id>")
async def progress_stream(workflow_id):
    async def generate():
        handle = client.get_workflow_handle(workflow_id)

        while True:
            describe_response = await handle.describe()
            if describe_response.status == WorkflowExecutionStatus.RUNNING:
                yield f"data: running\n\n"
                await asyncio.sleep(1)
            else:
                yield f"data: {describe_response.status.name.lower()}\n\n"
                break

    return Response(generate(), mimetype='text/event-stream')


@app.route("/book_result/<workflow_id>")
async def book_result(workflow_id):
    # Get workflow handle and wait for result
    handle = client.get_workflow_handle(workflow_id)
    result = await handle.result()

    if result == "Booking cancelled":
        return await render_template("book_result.html", cancelled=True)

    else:
        print(f"Result: {result}")
        result_list = result.split("Booked ")
        flight = result_list[1].split(": ")[1].title()
        hotel = result_list[2].split(": ")[1].title()
        car = result_list[3].split(": ")[1].title()

        return await render_template(
            "book_result.html",
            result=result,
            cancelled=False,
            flight=flight,
            hotel=hotel,
            car=car,
            user_id=workflow_id,
        )


if __name__ == "__main__":
    app.run(host="127.0.0.1", debug=True)
