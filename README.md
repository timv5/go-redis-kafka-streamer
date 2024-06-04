# go-redis-kafka-streamer
Project made for learning purpose of golang & redis

## Description
This is a simple backend application containing REST api for saving data
into Redis cache and then retrieving it - it's also saved into db.
There are 2 endpoints: 

1.
POST: http://localhost:8080/api/message/send

request:
{
"header": "header example",
"body": "body example"
}

response:
{
"ID": "831fee83-dcdc-59f2-87f9-c3d23e00b62d",
"CreatedAt": "2024-06-01T11:39:31.275707+02:00",
"Header": "header example",
"Body": "body example"
}

2.
POST: http://localhost:8080/api/message/retrieve

request:
{
"uuid": "87a1a6a7-e8d1-5c39-a562-e0b523382cf3"
}

response:
{
"ID": "831fee83-dcdc-59f2-87f9-c3d23e00b62d",
"Header": "header example",
"Body": "body example"
}

## Used technologies
- golang
- docker, docker compose
- redis