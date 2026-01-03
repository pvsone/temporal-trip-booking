require 'json/add/struct'

module TripBooking
  module Models
    BookTripInput = Struct.new(:userId, :flightId, :hotelId, :carId) do
    end
  end
end
