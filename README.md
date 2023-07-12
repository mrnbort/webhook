# Webhook Service 

## Description

The Webhook Service allows the user to execute predefined 
tasks based on the task id and the authentication key. The service 
has a REST-like interface with a single entry point.

## Technical details

The service loads a YAML configuration file which contains the list 
of tasks, where each task is associated with a command
that needs to be executed and its authentication key.
The REST-like API has a single entry point: POST. A POST request 
accepts a task id and an authentication key. After a POST request is
passed to the service, it locates the corresponding task
in the configuration file and executes the related command 
only if the authentication key matches.

## Run in Docker

1. Copy docker-compose.yml

    - change the ports if needed

2. Start a container with `docker-compose up`

## API

### Public Endpoints

`POST /execute?id={task_id}&key={key}` - executes the command 
for the requested task ID if the authentication key matches

- Returns:
    ```json
    {
    "status": "ok" 
    }
    ```

## command line parameters

```
Application Options:
     --port=            http data server port (default: 8080)
     --file             configuration file name (default: config.yaml)
	
Help Options:
 -h, --help                Show this help message
```
