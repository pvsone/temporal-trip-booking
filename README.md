# Temporal Trip Booking

A demonstration of Temporal workflows implementing a trip booking system that orchestrates flight, hotel, and car reservations with automatic compensation using the Saga pattern.

## Overview

This project showcases how to build reliable, distributed applications using Temporal. It implements a trip booking workflow that:
- Books a flight, hotel, and car in sequence
- Handles failures gracefully with automatic compensation
- Demonstrates retry strategies for transient failures
- Uses the Saga pattern to ensure transactional consistency

## Features

- **Workflow Orchestration**: Coordinates multiple booking activities in a reliable, durable manner
- **Saga Pattern**: Implements compensation logic to undo bookings if any step fails
- **Error Handling**: Demonstrates different error types:
  - **Transient errors** (retryable): Simulated with "flaky" flight IDs
  - **Application errors** (non-retryable): Simulated with "invalid" car IDs
  - **Logical bugs**: Simulated with "buggy" hotel IDs
- **Automatic Retries**: Configurable retry policies with exponential backoff
- **Testing**: Includes unit tests using Temporal's test framework

## Running

### Prerequisites
Install the language tools using [asdf](https://asdf-vm.com/guide/getting-started.html)

```bash
# Add the plugins
asdf plugin add dotnet
asdf plugin add golang
asdf plugin add java
asdf plugin add nodejs
asdf plugin add python
asdf plugin add ruby
asdf plugin add uv

# Install the tools in the .tool-versions file
asdf install
```

Install the [Temporal CLI](https://docs.temporal.io/cli)

### Start the Temporal server
```bash
temporal server start
```

### Run the UI
```bash
cd ui
./start_ui.sh
```

### Run the worker

Set the profile for the Temporal server connection
- (Optional) `export TEMPORAL_PROFILE=default`

Choose one language runtime:
- .NET: `cd dotnet && ./start_worker.sh`
- Go: `cd go && ./start_worker.sh`
- Java: `cd java && ./start_worker.sh`
- Python: `cd python && ./start_worker.sh`
- Ruby: `cd python && ./start_worker.sh`
- TypeScript: `cd typescript && ./start_worker.sh`

## Demo

### Happy Path
Enter your booking information in the Flask app <http://127.0.0.1:5000>, then see the tasks in the Web UI at
<http://localhost:8233/>.

### Durable Path
Use ctrl+c to kill the worker process at any time during workflow execution.  Restart the worker to resume the
function execution in a new process.

### Recover Forward (retries)
Add the word "Flaky" to the Flight input field.  The flight booking activity will fail 5 times, and succeed on the
6th attempt.

### Recover Forward (bug in code)
Add the word "Buggy" to the Hotel input field.  The hotel booking activity will fail indefinitely - until the bug in
the activity function is fixed, and the worker is restarted.

### Recover Backward (rollback)
Add the word "Invalid" to the Car input field.  The car booking activity will fail with and non-retryable error,
and the flight and hotel bookings will be rolled back.
