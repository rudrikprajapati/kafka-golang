# Kafka Report Generator Simulation

This project simulates 500 users uploading data to generate reports using Kafka, with a producer (Gin server) and consumer (worker) architecture. It runs locally via Docker Compose and tests two scenarios: 1 worker (~83 minutes) and 10 workers (~8.3 minutes).

## Project Structure

```
kafka/
├── consumer/
│   ├── consumer.go
│   └── Dockerfile
├── producer/
│   ├── producer.go
│   └── Dockerfile
├── go.mod
├── go.sum
├── docker-compose.yml
└── simulate.sh
```

## Prerequisites

- Docker and Docker Compose installed.
- Go 1.23+ (for local builds, optional).

## Setup

1. **Clone the Repository** (if applicable):

   ```bash
   git clone <repository-url>
   cd kafka
   ```

2. **Ensure Files**:
   - `producer.go`, `consumer.go`, and Dockerfiles use `producer.go`/`consumer.go` names.
   - `docker-compose.yml` includes `init-kafka` with 10 partitions active (1 partition commented out).

## Running the Simulation

### Scenario 1: Single Worker (~83 Minutes)

Tests 500 tasks processed by 1 consumer.

1. **Edit `docker-compose.yml`**:

   - Comment out the 10-partition `init-kafka` command:
     ```yaml
     # init-kafka:
     #   image: bitnami/kafka:latest
     #   command: >
     #     sh -c "until kafka-topics.sh --list --bootstrap-server kafka:9092; do echo 'Waiting for Kafka...'; sleep 2; done && kafka-topics.sh --create --topic report-tasks --bootstrap-server kafka:9092 --partitions 10 --replication-factor 1 || true"
     ```
   - Uncomment the 1-partition version:
     ```yaml
     init-kafka:
       image: bitnami/kafka:latest
       command: >
         sh -c "until kafka-topics.sh --list --bootstrap-server kafka:9092; do echo 'Waiting for Kafka...'; sleep 2; done && kafka-topics.sh --create --topic report-tasks --bootstrap-server kafka:9092 --partitions 1 --replication-factor 1 || true"
       depends_on:
         kafka:
           condition: service_healthy
     ```

2. **Start the Stack**:

   ```bash
   docker-compose up --build
   ```

3. **Simulate 500 Users**:

   ```bash
   ./simulate.sh 500
   ```

   - Expect: `<1s` total send time.

4. **Monitor Consumer**:

   ```bash
   docker-compose logs consumer --follow
   ```

   - Expect: 1 consumer processing 500 tasks (~83min).
   - Stop early if desired:
     ```bash
     docker-compose down
     ```

5. **Validate**:
   - Producer: <1s (from `simulate.sh`), ~50-100ms per request (logs).
   - Consumer: ~83min (extrapolate from logs).

### Scenario 2: 10 Workers (~8.3 Minutes)

Tests 500 tasks processed by 10 consumers.

1. **Update init-kafka**:

   - Modify docker-compose.yml to create the topic with 10 partitions initially, avoiding future override:
     ```yaml
     init-kafka:
       image: bitnami/kafka:latest
       command: >
         sh -c "until kafka-topics.sh --list --bootstrap-server kafka:9092; do echo 'Waiting for Kafka...'; sleep 2; done && kafka-topics.sh --create --topic report-tasks --bootstrap-server kafka:9092 --partitions 10 --replication-factor 1 || true"
       depends_on:
         kafka:
           condition: service_healthy
     ```

2. **Start Kafka and ZooKeeper**:

   ```bash
   docker-compose up -d kafka zookeeper
   ```

3. **Manually Create or Alter Topic**:

   - **If Topic Doesn't Exist**:
     ```bash
     docker exec -it kafka-report-sim-kafka-1 kafka-topics.sh --create --topic report-tasks --bootstrap-server kafka:9092 --partitions 10 --replication-factor 1
     ```
   - **If Topic Exists**:
     ```bash
     docker exec -it kafka-report-sim-kafka-1 kafka-topics.sh --alter --topic report-tasks --bootstrap-server kafka:9092 --partitions 10
     ```
   - **Verify**:
     ```bash
     docker exec -it kafka-report-sim-kafka-1 kafka-topics.sh --describe --topic report-tasks --bootstrap-server kafka:9092
     ```
     - Confirm PartitionCount: 10.

4. **Run Full Stack with 10 Consumers**:

   ```bash
   docker-compose up --build --scale consumer=10
   ```

5. **Simulate 500 Users** :

   ```bash
   ./simulate.sh 500
   ```

   - Expect: `<1s` total send time.

6. **Check Consumer Logs**:

   ```bash
   docker-compose logs consumer --follow
   ```

   - **Expected**: All 10 consumers (e.g., consumer-1 to consumer-10) should log "Starting report for user...", with tasks distributed roughly evenly (~50 tasks each).
   - **Timing**: Total ~8.3 minutes (500 tasks / 10 workers = 50 tasks × 10s).
   - Stop early if desired:
     ```bash
     docker-compose down
     ```

7. **Validate**:
   - Producer: <1s (from `simulate.sh`), ~50-100ms per request (logs).
   - Consumer: ~8.3min (all 10 active, extrapolate if stopped).

## Expected Results

- **Single Worker**: Producer <1s, Consumer ~83min.
- **10 Workers**: Producer <1s, Consumer ~8.3min.

## Troubleshooting

- **Producer Slow**: Check `docker-compose logs producer`.
- **Consumers Inactive**: Verify partitions:
  ```bash
  docker exec -it kafka-kafka-1 kafka-topics.sh --describe --topic report-tasks --bootstrap-server kafka:9092
  ```
- **Errors** :`docker-compose logs <service>`.

## Notes

- `simulate.sh` sends 500 requests concurrently.
- Consumer processing simulates a 10s report generation delay.
- Adjust `init-kafka` in `docker-compose.yml` to switch scenarios.
