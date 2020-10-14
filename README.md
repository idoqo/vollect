# Vollect
A naive job queue (and implementations) with support for pause/resume/terminate.

## Running the project
Build and run the project with:
```bash
$ docker-compose up --build
```
The project includes two implementations of the queue:

- `counter.go`: This implements a simple counter. Its job is to print a sequence of numbers (0-1500 if you let it) 
to stdout. It also sleeps for 3 seconds to simulate a delayed job. To create a counter task, make a HTTP 
GET request to `http://localhost:8080/counter`.
- `csv.go`: This implements a basic CSV parser by printing each row of the CSV file to stdout. To create
a CSV task, make a HTTP POST request to `http://localhost:8080/csv/upload`. The request should contain a file field
named `csv`. Here's an example that uses [httpie](httpie.org) and the `addresses.csv` file included in the repo.
```
$ http --form localhost:8080/csv/upload csv@addresses.csv
```
You should get a response that includes the ID of the created task e.g:
```json
{
    "message": "file upload complete. Task ID is 1"
}
```

## Endpoints
### Get All Tasks
This returns a list of all the jobs in the queue (both pending, running and paused).
#### Url
`/tasks`
#### Response

### Pause a Task
Pauses a given task using its ID. 
#### Url
`POST /tasks/{taskId}/pause` e.g `httpie POST localhost:8080/tasks/1/pause`
#### 
#### Response
The response is a 200 OK header with an empty body.

### Resume a Task
Resumes a paused task using its ID
#### Url
`POST /tasks/{taskId}/resume` e.g `httpie POST localhost:8080/tasks/1/resume`
#### Response
The response is a 200 OK header with an empty body.

### Terminates a Task
#### Url
`POST /tasks/{taskId}/terminate` e.g `httpie POST localhost:8080/tasks/1/terminate`
#### Response
The response is a 200 OK header with an empty body.

## How it works
The _queue_ contains *Tasks* which are saved to a PostgreSQL database. The tasks are run sequentially
 that is, only one task is running at a given time. A task has:

- ID: The task ID in the database
- Name: A random name for the task, should be unique.
- Payload: Payload contains metadata about the task such as the "handler" that should run the task.
The payload can also be updated before running task, that way, the task can resume from that state on resumption.
- Handler: Handlers should implement all methods of the `TaskHandler` interface 
(i.e., `Handle`, `OnPause`, and `OnResume`). 

Each queue implementation creates its own `Handle` method specifying how the task should be executed.
The Handle method takes in an integer channel as parameter.
This channel listens for signal to **pause** - so that we can pause the task even if it's already running.

`OnPause` implementations should return the current state of the task, represented as `map[string]interface{}`.

`OnResume` takes in a task state and uses it to restore the task to where it was before it was paused.

Terminated tasks get removed from the database entirely.

## Some Issues

- Couldn't figure how to jump to specific row while reading CSVs (using `encoding/csv` package).




