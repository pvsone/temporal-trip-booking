# frozen_string_literal: true

require_relative 'shared'
require 'temporalio/activity'

module TripBooking
  module Activities
    class BookFlight < Temporalio::Activity::Definition
      def execute(input)
        Temporalio::Activity::Context.current.logger.info("Booking flight: #{input.flightId}")
        sleep(1)

        if input.flightId.to_s.downcase.include?('flaky')
          # a transient error, which will be retried
          attempt = Temporalio::Activity::Context.current.info.attempt
          raise 'Flight booking service is currently unavailable' if attempt < 6
        end

        "Booked flight: #{input.flightId}"
      end
    end

    class BookHotel < Temporalio::Activity::Definition
      def execute(input)
        Temporalio::Activity::Context.current.logger.info("Booking hotel: #{input.hotelId}")
        sleep(1)

        if input.hotelId.to_s.downcase.include?('buggy')
          # a simulated bug
          error = true
          raise 'Error due to bug in code' if error
        end

        "Booked hotel: #{input.hotelId}"
      end
    end

    class BookCar < Temporalio::Activity::Definition
      def execute(input)
        Temporalio::Activity::Context.current.logger.info("Booking car: #{input.carId}")
        sleep(1)

        if input.carId.to_s.downcase.include?('invalid')
          # a business error, which cannot be retried
          raise Temporalio::Error::ApplicationError.new(
            "Car #{input.carId} is invalid",
            non_retryable: true
          )
        end

        "Booked car: #{input.carId}"
      end
    end

    class NotifyUser < Temporalio::Activity::Definition
      def execute(input)
        Temporalio::Activity::Context.current.logger.info("Notifying user: #{input.userId}")
        sleep(1)
        "Notified user: #{input.userId}"
      end
    end

    class UndoBookFlight < Temporalio::Activity::Definition
      def execute(input)
        Temporalio::Activity::Context.current.logger.info("Undo flight booking: #{input.flightId}")
        sleep(1)
        "Unbooked flight: #{input.flightId}"
      end
    end

    class UndoBookHotel < Temporalio::Activity::Definition
      def execute(input)
        Temporalio::Activity::Context.current.logger.info("Undo hotel booking: #{input.hotelId}")
        sleep(1)
        "Unbooked hotel: #{input.hotelId}"
      end
    end

    class UndoBookCar < Temporalio::Activity::Definition
      def execute(input)
        Temporalio::Activity::Context.current.logger.info("Undo car booking: #{input.carId}")
        sleep(1)
        "Unbooked car: #{input.carId}"
      end
    end
  end
end
