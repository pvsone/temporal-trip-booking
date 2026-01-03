package app

type BookTripInput struct {
	UserId   string `json:"userId"`
	FlightId string `json:"flightId"`
	HotelId  string `json:"hotelId"`
	CarId    string `json:"carId"`
}
