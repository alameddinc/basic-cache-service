## How is Works
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
