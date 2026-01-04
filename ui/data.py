from dataclasses import dataclass


@dataclass
class BookTripInput:
    userId: str
    flightId: str
    hotelId: str
    carId: str
