import { ApplicationFailure, Context, log, sleep } from '@temporalio/activity';
import type { BookTripInput } from './shared';

export async function bookFlight(input: BookTripInput): Promise<string> {
  log.info(`Booking flight: ${input.flightId}`);

  await sleep(1000);

  if (input.flightId.toLowerCase().includes('flaky')) {
    // a transient error, which will be retried
    if (Context.current().info.attempt < 6) {
      throw ApplicationFailure.create({
        message: 'Flight booking service is currently unavailable',
      });
    }
  }

  return `Booked flight: ${input.flightId}`;
}

export async function bookHotel(input: BookTripInput): Promise<string> {
  log.info(`Booking hotel: ${input.hotelId}`);

  await sleep(1000);

  if (input.hotelId.toLowerCase().includes('buggy')) {
    // a simulated bug
    const error = true;
    if (error) {
      throw new Error('Error due to bug in code');
    }
  }

  return `Booked hotel: ${input.hotelId}`;
}

export async function bookCar(input: BookTripInput): Promise<string> {
  log.info(`Booking car ${input.carId}`);

  await sleep(1000);

  if (input.carId.toLowerCase().includes('invalid')) {
    // a business error, which cannot be retried
    throw ApplicationFailure.create({
      nonRetryable: true,
      message: `Car ${input.carId} is invalid`,
      type: 'InvalidCar',
    });
  }

  return `Booked car: ${input.carId}`;
}

export async function notifyUser(input: BookTripInput): Promise<string> {
  log.info(`Notifying user: ${input.userId}`);

  await sleep(1000);

  return `Notified user: ${input.userId}`;
}

export async function undoBookFlight(input: BookTripInput): Promise<string> {
  log.info(`Undo flight booking: ${input.flightId}`);

  await sleep(1000);

  return `Unbooked flight: ${input.flightId}`;
}

export async function undoBookHotel(input: BookTripInput): Promise<string> {
  log.info(`Undo hotel booking: ${input.hotelId}`);

  await sleep(1000);

  return `Unbooked hotel: ${input.hotelId}`;
}

export async function undoBookCar(input: BookTripInput): Promise<string> {
  log.info(`Undo car booking: ${input.carId}`);

  await sleep(1000);

  return `Unbooked car: ${input.carId}`;
}
