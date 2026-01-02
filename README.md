# Temporal Trip Booking

## Overview

The Temporal trip booking workflow is responsible for coordinating the booking of a vacation package, consisting of a
flight, hotel, and car reservation. In the event of a failure at any point during the booking process, the workflow
will trigger compensating actions to undo any previous bookings. This is an implementation of a Saga pattern. A Saga
ensures that a business process is atomic, that is, it executes observably equivalent to completely or not at all.

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
Choose one language runtime:

- Python: `cd python && ./start_worker.sh`
- TypeScript: `cd typescript && ./start_worker.sh`
- Go: `cd go && ./start_worker.sh`
- Java: `cd java && ./start_worker.sh`
- .NET: `cd dotnet && ./start_worker.sh`

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
