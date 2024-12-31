# Temporal Trip Booking

## Overview

The Temporal trip booking workflow is responsible for coordinating the booking of a vacation package, consisting of a car, hotel, and flight
reservation. In the event of a failure at any point during the booking process, the workflow will trigger compensating actions to undo any
previous bookings. This is an implementation of a Saga pattern. A Saga ensures that a business process is atomic, that is, it executes observably
equivalent to completely or not at all.

## Running

Prerequisites:

- Python >= 3.7
- [Poetry](https://python-poetry.org)
- [Local Temporal server running](https://docs.temporal.io/application-development/foundations#run-a-development-cluster)

With this repository cloned, run the following at the root of the directory:

```bash
poetry install --no-root
```
That loads all required dependencies.

Then run the worker and workflow.

```bash
cd app-python
./start_worker.sh

cd ui
./start_ui.sh
```

### Demo: Happy Path
Enter your booking information in the Flask app <http://127.0.0.1:5000>, then see the tasks in the Web UI at <http://localhost:8233/>.

### Demo: Durable Path
Use ctrl+c to kill the worker process at any time during workflow execution.  Restart the worker to resume the function execution in a new process.

### Demo: Recover Forward (retries)
Add the word "Flaky" to the Car input field.  The car booking activity will fail 5 times, and succeed on the 6th attempt.

### Demo: Recover Forward (bug in code)
Add the word "Buggy" to the Hotel input field.  The hotel booking activity will fail indefinitely - until the bug in the activity function is fixed, and the worker is restarted.

### Demo: Recover Backward (rollback)
Add the word "Invalid" to the Flight input field.  The flight booking activity will fail with and non-retryable error, and the hotel and car bookings will be rolled back.
