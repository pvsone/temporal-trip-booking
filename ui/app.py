from flask import Flask, render_template, request
import uuid
from temporalio.client import Client
from dataclasses import dataclass


@dataclass
class BookTripInput:
    userId: str
    flightId: str
    hotelId: str
    carId: str


app = Flask(__name__)


@app.route("/")
async def display_form():
    return render_template("book_trip.html")


@app.route("/book", methods=["POST"])
async def book_trip():
    user_id = f'{request.form.get("name").replace(
        " ", "-").lower()}-{str(uuid.uuid4().int)[:6]}'
    flight = request.form.get("flight")
    hotel = request.form.get("hotel")
    car = request.form.get("car")

    input = BookTripInput(
        userId=user_id,
        flightId=flight,
        hotelId=hotel,
        carId=car,
    )

    client = await Client.connect("localhost:7233")

    result = await client.execute_workflow(
        "BookWorkflow",
        input,
        id=user_id,
        task_queue="trip-task-queue",
    )
    if result == "Booking cancelled":
        return render_template("book_trip.html", cancelled=True)

    else:
        print(result)
        result_list = result.split("Booked ")
        flight = result_list[1].split(": ")[1].title()
        hotel = result_list[2].split(": ")[1].title()
        car = result_list[3].split(": ")[1].title()

        print(user_id)
        return render_template(
            "book_trip.html",
            result=result,
            cancelled=False,
            flight=flight,
            hotel=hotel,
            car=car,
            user_id=user_id,
        )


if __name__ == "__main__":
    app.run(host="127.0.0.1", debug=True)
