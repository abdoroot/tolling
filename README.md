# Tolling

`tolling` is a small Go microservice demo for collecting vehicle location events, calculating traveled distance, and turning that into a toll invoice.

The system uses:

- WebSockets for ingesting OBU telemetry
- Kafka for event transport between services
- gRPC and HTTP for internal service communication
- Prometheus metrics on the aggregator

## Architecture

The runtime flow in this repository is:

1. `OBU` opens a WebSocket connection to the data receiver and sends random latitude/longitude pairs for 20 generated OBU IDs every 5 seconds.
2. `data_receiver` accepts the WebSocket payloads on port `3000` and publishes them to the Kafka topic `data`.
3. `distance_calculator` consumes the `data` topic, calculates a distance delta, and forwards it to the aggregator over gRPC at `127.0.0.1:3002`.
4. `aggregator` stores total distance in memory, exposes invoice lookup over HTTP on `3001`, and publishes Prometheus metrics on `2112`.
5. `gateway` exposes a simpler public HTTP endpoint on `6000` that proxies invoice reads to the aggregator.

## Services

| Component | Entry point | Protocol | Default address | Responsibility |
| --- | --- | --- | --- | --- |
| OBU simulator | `OBU/main.go` | WebSocket client | `ws://localhost:3000` | Generates random OBU position data |
| Data receiver | `data_receiver/main.go` | HTTP/WebSocket server | `:3000` | Receives OBU payloads and produces Kafka messages |
| Distance calculator | `distance_calculator/main.go` | Kafka consumer + gRPC client | Kafka `localhost`, gRPC `127.0.0.1:3002` | Converts telemetry into distance events |
| Aggregator | `aggregator/main.go` | HTTP server + gRPC server | HTTP `:3001`, gRPC `:3002`, metrics `:2112` | Aggregates distance and calculates invoice totals |
| Gateway | `gateway/main.go` | HTTP server | `:6000` | Public invoice lookup endpoint |

## Repository Layout

```text
.
├── OBU/                    # OBU telemetry simulator
├── aggregator/             # Distance aggregation, invoice calculation, HTTP/gRPC transports
├── config/prometheus.yml   # Example Prometheus scrape config
├── data_receiver/          # WebSocket ingestion and Kafka producer
├── distance_calculator/    # Kafka consumer and gRPC aggregation client
├── gateway/                # Public invoice API
├── types/                  # Shared types and protobuf definitions
├── docker-compose.yml      # Kafka + Zookeeper
└── Makefile                # Convenience run targets
```

## Prerequisites

- Go `1.18+`
- Docker and Docker Compose
- `protoc` only if you need to regenerate protobuf files
- A local Kafka-compatible runtime reachable as `localhost` (the provided Compose file exposes Kafka on `localhost:9092`)

The Kafka client library in this repo is `confluent-kafka-go`, which usually requires `librdkafka` development libraries to be installed on the host when building locally.

## Run Locally

Start Kafka and Zookeeper:

```bash
docker compose up -d
```

Run the services in separate terminals from the repository root:

```bash
make agg
make receiver
make calc
make gate
make obu
```

Recommended startup order:

1. `make agg`
2. `make receiver`
3. `make calc`
4. `make gate`
5. `make obu`

Once everything is running, the OBU simulator will continuously produce location updates, and invoices will begin accumulating in the aggregator's in-memory store.

## Available Commands

```bash
make obu       # run the OBU simulator
make receiver  # run the data receiver
make calc      # run the distance calculator
make agg       # run the aggregator
make gate      # run the gateway
make proto     # regenerate protobuf and gRPC Go files
```

## HTTP and gRPC Interfaces

### Gateway API

Get an invoice through the gateway:

```bash
curl "http://127.0.0.1:6000/invoice?obuid=123"
```

### Aggregator HTTP API

Get an invoice directly from the aggregator:

```bash
curl "http://127.0.0.1:3001/invoice?obuid=123"
```

Manually aggregate a distance event:

```bash
curl -X POST "http://127.0.0.1:3001/aggregate" \
  -H "Content-Type: application/json" \
  -d '{"obuid":123,"value":42.5,"unix":1710000000}'
```

The invoice response shape is:

```json
{
  "obuid": 123,
  "total_distance": 42.5,
  "total_amount": 96.05
}
```

`total_amount` is calculated as:

```text
total_distance * 2.26
```

### Aggregator gRPC API

The protobuf contract lives in [`types/ptypes.proto`](./types/ptypes.proto).

The `distance_calculator` currently uses the gRPC method below to submit new distance events:

- `Aggreagator.AggregateDistance(DistanceRequest) returns (None)`

## Metrics

The aggregator exposes Prometheus metrics at:

```text
http://127.0.0.1:2112/metrics
```

This repo includes an example Prometheus scrape configuration at `config/prometheus.yml`.

If you have a local Prometheus binary, you can start it with:

```bash
./prometheus --config.file=./config/prometheus.yml
```

## Current Implementation Notes

These details are important if you are using this repository as a reference or starting point:

- Aggregated distance is stored only in memory. Restarting the aggregator clears all totals.
- The distance calculator keeps a single previous point for the whole process, not one previous point per OBU. In the current implementation, consecutive messages from different vehicles can affect each other's calculated distance.
- The OBU simulator generates random coordinates and random OBU IDs on startup, so values are not deterministic across runs.
- There are no automated tests in the repository at the moment.
- There is no authentication, persistence layer, or deployment configuration beyond the local Kafka Compose file.

## Shared Data Types

The main shared payloads are:

- `types.OBUdata`: `{ obuid, lat, long }`
- `types.Distance`: `{ obuid, value, unix }`
- `types.Invoice`: `{ obuid, total_distance, total_amount }`

## Development Notes

- Kafka topic name is hard-coded as `data`.
- The gateway reads from the aggregator over HTTP.
- The distance calculator sends writes to the aggregator over gRPC.
- Ports and endpoints are hard-coded in the source today; there is no environment-based configuration layer yet.
