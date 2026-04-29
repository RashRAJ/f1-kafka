## Areas of improvements
- increase the no of partitions, from current setting of 0
  - dont want to create topic from producer (helps scalablilty without code changes)
  -  why is increaing partition increment beneficial? {Scalability and redundancy}
  - best to make at infra level
  - Topic replicas
- access realtime data without hitting the lilmitation of fetching too much data at once
   - implement ratelimiting on api endpoint
   - properly configure http ratelimit on request 
   - event driven
   - test the usage of REST API, MQTT, and/or WebSockets. for realtime consumption
- send data securely 
- http/Rest
  - centralize your confg, move base url
  - set up an http client, with a context timeout ref {https://github.com/Harphies/go.microservices.io/blob/main/utils/http-helpers-new.go}
  - implement buffered reading 
- producer telemetry {data, duration, error rate, transfer rate, etc set SLO}
- retry and error handling
- DLQ
- reconfigure producer, producer should run as a thread
- create a producer config
- confluent cloud kafka.


#################
# focused learning

- Kafka auth
- kafka observability
- mtls auth
- API Design








Main could still be improved


2025/12/22 16:08:31 Failed to fetch team radio: API returned status 429: {"detail":"Data limit exceeded. Max 4MB per 10 second(s). Try again in 2 seconds.","error":"Too Many Requests"}

2025/12/22 16:11:05 Failed to send intervals to Kafka: kafka: invalid configuration (Attempt to produce message larger than configured Producer.MaxMessageBytes: 2877851 > 1048576)