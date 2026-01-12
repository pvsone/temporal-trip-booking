package app

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.temporal.io/sdk/testsuite"
)

type UnitTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite

	env *testsuite.TestWorkflowEnvironment
}

func (s *UnitTestSuite) SetupTest() {
	s.env = s.NewTestWorkflowEnvironment()
}

func (s *UnitTestSuite) AfterTest(suiteName, testName string) {
	s.env.AssertExpectations(s.T())
}

func TestUnitTestSuite(t *testing.T) {
	suite.Run(t, new(UnitTestSuite))
}

func (s *UnitTestSuite) Test_BookWorkflow_With_Activities() {

	bookTripInput := BookTripInput{
		UserId:   "Test User",
		FlightId: "Flight 123",
		HotelId:  "Hotel 456",
		CarId:    "Car 789",
	}

	s.env.RegisterActivity(BookFlight)
	s.env.RegisterActivity(BookHotel)
	s.env.RegisterActivity(BookCar)
	s.env.RegisterActivity(NotifyUser)
	// s.env.RegisterActivity(UndoBookFlight)
	// s.env.RegisterActivity(UndoBookHotel)
	// s.env.RegisterActivity(UndoBookCar)

	s.env.ExecuteWorkflow(BookWorkflow, bookTripInput)
	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())

	var result string
	s.NoError(s.env.GetWorkflowResult(&result))
	s.Equal("Booked flight: Flight 123 Booked hotel: Hotel 456 Booked car: Car 789", result)
}

func (s *UnitTestSuite) Test_BookWorkflow_With_Activities_Mocked() {

	bookTripInput := BookTripInput{
		UserId:   "Test User",
		FlightId: "Flight 123",
		HotelId:  "Hotel 456",
		CarId:    "Car 789",
	}

	s.env.OnActivity(BookFlight, mock.Anything, mock.Anything).Return("Booked flight: Flight 123", nil)
	s.env.OnActivity(BookHotel, mock.Anything, mock.Anything).Return("Booked hotel: Hotel 456", nil)
	s.env.OnActivity(BookCar, mock.Anything, mock.Anything).Return("Booked car: Car 789", nil)
	s.env.OnActivity(NotifyUser, mock.Anything, mock.Anything).Return("Notified user: Test User", nil)

	s.env.ExecuteWorkflow(BookWorkflow, bookTripInput)
	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())

	var result string
	s.NoError(s.env.GetWorkflowResult(&result))
	s.Equal("Booked flight: Flight 123 Booked hotel: Hotel 456 Booked car: Car 789", result)
}
