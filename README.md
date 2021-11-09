## What is It?
It is a simple caching system written in golang. This saves the new values to new storage files with an interval time loop and if the values need to be changed or deleted, it performs operations on the recorded file with the multi-process feature. Every storagefile have multi value and when storage file have never value, It will be deleted.

## How It Runs?
With Docker-compose:
`docker-compose up --build -d`

Without Docker-compose:
`go run app\main.go`

## Endpoints
 #### Get Endpoint
 * Method: GET
 * Url : localhost:8080/storage/get/{key}
 * Body: `-`
 * Response: 
```
{"value":"value","storage":"stored_filename"}
```

#### Set Endpoint
 * Method: POST
 * Url : localhost:8080/storage/set
 * Body:
```
{
	"key":"key",
	"value":"value"
}
```
 * Response: 
```
{
	"value":"value",
	"storage":"stored_filename"
}
```

#### Flush Endpoint
 * Method: GET
 * Url : localhost:8080/storage/flush
 * Body: `-`
 * Response: 
```
{"message":"ok"}
```
