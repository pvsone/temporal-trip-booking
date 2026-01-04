import { ActivityFailure, ApplicationFailure, log, proxyActivities, sleep } from '@temporalio/workflow';
import type * as activities from './activities';
import type { BookTripInput } from './shared';

const { bookFlight, bookHotel, bookCar, notifyUser, undoBookFlight, undoBookHotel, undoBookCar } = proxyActivities<
  typeof activities
>({
  startToCloseTimeout: '5s',
  retry: {
    initialInterval: '1s',
    backoffCoefficient: 2,
    maximumInterval: '30s',
  },
});

interface Compensation {
  message: string;
  fn: () => Promise<void>;
}

export async function BookWorkflow(input: BookTripInput): Promise<string> {
  log.info(`Book workflow started, user_id = ${input.userId}`);

  // Saga compensations
  const compensations: Compensation[] = [];

  try {
    // Book Flight
    const flight = await bookFlight(input);
    compensations.unshift({
      message: prettyErrorMessage('undo flight booking'),
      fn: async () => {
        await undoBookFlight(input);
      },
    });

    log.info('Sleeping for 1 second...');
    await sleep('1 seconds');

    // Book Hotel
    const hotel = await bookHotel(input);
    compensations.unshift({
      message: prettyErrorMessage('undo hotel booking'),
      fn: async () => {
        await undoBookHotel(input);
      },
    });

    log.info('Sleeping for 1 second...');
    await sleep('1 seconds');

    // Book Car
    const car = await bookCar(input);
    compensations.unshift({
      message: prettyErrorMessage('undo car booking'),
      fn: async () => {
        await undoBookCar(input);
      },
    });

    // Send Notification
    await notifyUser(input);

    return `${flight} ${hotel} ${car}`;
  } catch (err) {
    if (err instanceof ActivityFailure && err.cause instanceof ApplicationFailure) {
      log.error(err.cause.message);
    } else {
      log.error(`Failed to book trip: ${err}`);
    }
    // an error occurred so call compensations
    await compensate(compensations);

    return 'Booking canceled';
  }
}

async function compensate(compensations: Compensation[] = []) {
  if (compensations.length > 0) {
    log.info('failures encountered during account opening - compensating');
    for (const comp of compensations) {
      try {
        log.error(comp.message);
        await comp.fn();
      } catch (err) {
        log.error(`failed to compensate: ${prettyErrorMessage('', err)}`, { err });
        // swallow errors
      }
    }
  }
}

function prettyErrorMessage(message: string, err?: any) {
  let errMessage = err && err.message ? err.message : '';
  if (err && err instanceof ActivityFailure) {
    errMessage = `${err.cause?.message}`;
  }
  return `${message}: ${errMessage}`;
}
