using System.Text.Json.Serialization;

namespace TemporalTripBooking;

public class Shared 
{
    [JsonPropertyName("userId")]
	public string UserId { get; set; }
    [JsonPropertyName("flightId")]
	public string FlightId { get; set; }
    [JsonPropertyName("hotelId")]
	public string HotelId { get; set; }
    [JsonPropertyName("carId")]
	public string CarId { get; set; }

	public Shared()
	{
		UserId = string.Empty;
		FlightId = string.Empty;
		HotelId = string.Empty;
		CarId = string.Empty;
	}

 	public Shared(string userId, string flightId, string hotelId, string carId)
	{
		UserId = userId;
		FlightId = flightId;
		HotelId = hotelId;
		CarId = carId;
	}
	
}
