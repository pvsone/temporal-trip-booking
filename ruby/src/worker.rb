# frozen_string_literal: true

require_relative 'activities'
require_relative 'workflow'
require 'logger'
require 'temporalio/client'
require 'temporalio/env_config'
require 'temporalio/worker'

puts "⚙️ Using TEMPORAL_PROFILE: '#{ENV['TEMPORAL_PROFILE']}'"
args, kwargs = Temporalio::EnvConfig::ClientConfig.load_client_connect_options
kwargs[:logger] = Logger.new($stdout, level: Logger::INFO)
client = Temporalio::Client.connect(*args, **kwargs)
puts "✅ Client connected to '#{args[0]}' in namespace '#{args[1]}'"

worker = Temporalio::Worker.new(
  client:,
  task_queue: 'trip-task-queue',
  workflows: [
    TripBooking::BookWorkflow
  ],
  activities: [
    TripBooking::Activities::BookFlight,
    TripBooking::Activities::BookHotel,
    TripBooking::Activities::BookCar,
    TripBooking::Activities::NotifyUser,
    TripBooking::Activities::UndoBookFlight,
    TripBooking::Activities::UndoBookHotel,
    TripBooking::Activities::UndoBookCar
  ]
)

# Run the worker until SIGINT. This can be done in many ways, see "Workers" section for details.
puts 'Starting Worker (press Ctrl+C to exit)'
worker.run(shutdown_signals: ['SIGINT'])
