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
      flight = Temporalio::Workflow.execute_activity(
        Activities::BookFlight,
        input,
        start_to_close_timeout: 5,
        retry_policy: retry_policy
      )
      compensations << Activities::UndoBookFlight

      Temporalio::Workflow.logger.info("Sleeping for 1 second...")
      Temporalio::Workflow.sleep(1)

      # Book hotel
      hotel = Temporalio::Workflow.execute_activity(
        Activities::BookHotel,
        input,
        start_to_close_timeout: 5,
        retry_policy: retry_policy
      )
      compensations << Activities::UndoBookHotel

      Temporalio::Workflow.logger.info("Sleeping for 1 second...")
      Temporalio::Workflow.sleep(1)

      # Book car
      car = Temporalio::Workflow.execute_activity(
        Activities::BookCar,
        input,
        start_to_close_timeout: 5,
        retry_policy: retry_policy
      )
      compensations << Activities::UndoBookCar

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
      "Booking canceled"
    end
  end
end
