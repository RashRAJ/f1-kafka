## Areas of improvements
- increase the no of partitions, from the current setting of 0
  - don't want to create topic from producer (helps scalability without code changes)
  -  Why is increasing the partition increment beneficial? {Scalability and redundancy}
  - best to make at the infra level
  - Topic replicas
- access real-time data without hitting the limitation of fetching too much data at once
   - implement ratelimiting on api endpoint
   - properly configure the HTTP rate limit on the request 
   - event-driven
   - test the usage of REST API, MQTT, and/or WebSockets, for real-time consumption
- send data securely 
- http/Rest
  - centralize your config, move base URL
  - set up an HTTP client, with a context timeout ref {https://github.com/Harphies/go.microservices.io/blob/main/utils/http-helpers-new.go}
  - implement buffered reading 
- producer telemetry {data, duration, error rate, transfer rate, etc., set SLO}
- retry and error handling
- DLQ
- reconfigure producer, producer should run as a thread
- create a producer config
- Confluent Cloud Kafka.


#################
# focused learning

- Kafka auth
- Kafka observability
- mtls auth
- API Design





