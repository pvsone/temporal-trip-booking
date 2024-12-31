import uuid

from flask import Flask, render_template, request
from temporalio.client import Client

from dataclasses import dataclass

@dataclass
class BookVacationInput:
    book_user_id: str
    book_car_id: str
    book_hotel_id: str
    book_flight_id: str

app = Flask(__name__)


@app.route("/")
async def display_form():
    return render_template("book_vacation.html")


@app.route("/book", methods=["POST"])
async def book_vacation():
    user_id = f'{request.form.get("name").replace(" ", "-").lower()}-{str(uuid.uuid4().int)[:6]}'
    car = request.form.get("car")
    hotel = request.form.get("hotel")
    flight = request.form.get("flight")

    input = BookVacationInput(
        book_user_id=user_id,
        book_car_id=car,
        book_hotel_id=hotel,
        book_flight_id=flight,
    )

    client = await Client.connect("localhost:7233")

    result = await client.execute_workflow(
        "BookWorkflow",
        input,
        id=user_id,
        task_queue="saga-task-queue",
    )
    if result == "Booking cancelled":
        return render_template("book_vacation.html", cancelled=True)

    else:
        print(result)
        result_list = result.split("Booked ")
        car = result_list[1].split(": ")[1].title()
        hotel = result_list[2].split(": ")[1].title()
        flight = result_list[3].split(": ")[1].title()
        print(user_id)
        return render_template(
            "book_vacation.html",
            result=result,
            cancelled=False,
            car=car,
            hotel=hotel,
            flight=flight,
            user_id=user_id,
        )


if __name__ == "__main__":
    app.run(host="127.0.0.1", debug=True)
