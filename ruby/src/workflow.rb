require_relative 'activities'
require_relative 'shared'
require 'temporalio/retry_policy'
require 'temporalio/workflow'

module TripBooking
  class BookWorkflow < Temporalio::Workflow::Definition
    def execute(input_hash)
      input = TripBooking::Models::BookTripInput.new(**input_hash)
      Temporalio::Workflow.logger.info("Book workflow started, userId = #{input.userId}")

      retry_policy = Temporalio::RetryPolicy.new(
        max_interval: 30
      )

      # Saga compensations
      compensations = []

      # Book flight
      compensations << Activities::UndoBookFlight
      flight = Temporalio::Workflow.execute_activity(
        Activities::BookFlight,
        input,
        start_to_close_timeout: 5,
        retry_policy: retry_policy
      )

      Temporalio::Workflow.logger.info("Sleeping for 1 second...")
      Temporalio::Workflow.sleep(1)

      # Book hotel
      compensations << Activities::UndoBookHotel
      hotel = Temporalio::Workflow.execute_activity(
        Activities::BookHotel,
        input,
        start_to_close_timeout: 5,
        retry_policy: retry_policy
      )

      Temporalio::Workflow.logger.info("Sleeping for 1 second...")
      Temporalio::Workflow.sleep(1)

      # Book car
      compensations << Activities::UndoBookCar
      car = Temporalio::Workflow.execute_activity(
        Activities::BookCar,
        input,
        start_to_close_timeout: 5,
        retry_policy: retry_policy
      )

      # Notify user
      Temporalio::Workflow.execute_activity(
        Activities::NotifyUser,
        input,
        start_to_close_timeout: 5,
        retry_policy: retry_policy
      )

      "#{flight}, #{hotel}, #{car}"
    rescue => error
      Temporalio::Workflow.logger.error("Failed to book trip #{error}")
      compensations.reverse_each do |compensation|
        Temporalio::Workflow.execute_activity(
          compensation,
          input,
          start_to_close_timeout: 5,
          retry_policy: retry_policy
        )
      end
      "Booking cancelled"
    end
  end
end