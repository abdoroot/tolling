# Tolling

`tolling` is a small Go microservice demo for collecting vehicle position events, calculating traveled distance, and turning that into a toll invoice.

The system uses:

- WebSockets for OBU telemetry ingestion
- Kafka for event transport
- gRPC for internal distance aggregation writes
- HTTP for invoice reads
- Prometheus for metrics

## Architecture

The intended topology is now split into edge-facing and internal services.

### Edge services

- `data_receiver` accepts WebSocket traffic from OBUs on port `3000`
- `gateway` exposes invoice lookup on port `6000`

### Internal services

- `aggregator` stores distance totals, calculates invoices, serves internal HTTP on `3001`, internal gRPC on `3002`, and metrics on `2112`
- `distance_calculator` consumes Kafka and forwards distance updates to the aggregator over gRPC
- Kafka and Zookeeper stay on the backend network for service-to-service traffic, with Kafka also published on `9092` for local development

## Runtime Flow

1. `OBU` opens a WebSocket connection to the data receiver and sends random latitude/longitude pairs for generated OBU IDs every 5 seconds.
2. `data_receiver` accepts the payloads on `:3000` and publishes them to the Kafka topic `data`.
3. `distance_calculator` consumes the `data` topic, calculates a distance delta, and forwards it to the aggregator over gRPC.
4. `aggregator` stores total distance in memory and calculates invoice totals.
5. `gateway` reads invoices from the aggregator's internal HTTP API and exposes the client-facing invoice endpoint.

## Services

| Component | Entry point | Protocol | Default address | Exposure |
| --- | --- | --- | --- | --- |
| OBU simulator | `OBU/main.go` | WebSocket client | `ws://localhost:3000` | client only |
| Data receiver | `data_receiver/main.go` | HTTP/WebSocket server | `:3000` | public |
| Distance calculator | `distance_calculator/main.go` | Kafka consumer + gRPC client | Kafka `localhost:9092`, gRPC `127.0.0.1:3002` | internal |
| Aggregator | `aggregator/main.go` | HTTP server + gRPC server + metrics | HTTP `127.0.0.1:3001`, gRPC `127.0.0.1:3002`, metrics `127.0.0.1:2112` | internal |
| Gateway | `gateway/main.go` | HTTP server | `:6000` | public |

Local defaults keep the aggregator bound to loopback so it is not exposed beyond the host by accident. In Docker Compose, the aggregator is attached only to the backend network and does not publish host ports.

## Repository Layout

```text
.
├── OBU/                         # OBU telemetry simulator
├── aggregator/                  # Distance aggregation, invoice calculation, HTTP/gRPC transports
├── config/prometheus.yml        # Prometheus config for local host runs
├── config/prometheus.docker.yml # Prometheus config for Docker Compose
├── data_receiver/               # WebSocket ingestion and Kafka producer
├── distance_calculator/         # Kafka consumer and gRPC aggregation client
├── gateway/                     # Public invoice API
├── internal/envutil/            # Small helpers for env-driven config
├── types/                       # Shared types and protobuf definitions
├── Dockerfile                   # Shared multi-stage image build for services
├── docker-compose.yml           # Full local stack
└── Makefile                     # Convenience run targets
```

## Configuration

The code no longer hard-codes all endpoints. These environment variables control the main runtime wiring:

| Variable | Default | Used by |
| --- | --- | --- |
| `OBU_WS_ENDPOINT` | `ws://localhost:3000` | `OBU` |
| `OBU_SEND_INTERVAL` | `5s` | `OBU` |
| `OBU_COUNT` | `20` | `OBU` |
| `DATA_RECEIVER_LISTEN_ADDR` | `:3000` | `data_receiver` |
| `KAFKA_BOOTSTRAP_SERVERS` | `localhost:9092` | `data_receiver`, `distance_calculator` |
| `KAFKA_TOPIC` | `data` | `data_receiver`, `distance_calculator` |
| `KAFKA_GROUP_ID` | `myGroup` | `distance_calculator` |
| `AGGREGATOR_HTTP_LISTEN_ADDR` | `127.0.0.1:3001` | `aggregator` |
| `AGGREGATOR_GRPC_LISTEN_ADDR` | `127.0.0.1:3002` | `aggregator` |
| `AGGREGATOR_METRICS_LISTEN_ADDR` | `127.0.0.1:2112` | `aggregator` |
| `AGGREGATOR_GRPC_ENDPOINT` | `127.0.0.1:3002` | `distance_calculator` |
| `GATEWAY_LISTEN_ADDR` | `:6000` | `gateway` |
| `GATEWAY_AGGREGATOR_ENDPOINT` | `http://127.0.0.1:3001/invoice` | `gateway` |

## Run Locally Without Docker

Prerequisites:

- Go `1.18+`
- Docker and Docker Compose for Kafka
- `protoc` only if you need to regenerate protobuf files
- `librdkafka` development libraries installed locally when building the Kafka-based services

Start Kafka and Zookeeper:

```bash
docker compose up -d zookeeper kafka
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

In this mode, the aggregator still exists as an internal service, but it is only bound to `127.0.0.1` by default.

## Run With Docker Compose

Build and start the core stack:

```bash
docker compose up --build
```

Start the optional OBU simulator too:

```bash
docker compose --profile simulator up --build
```

### Host-accessible endpoints in Compose

- WebSocket ingest: `ws://127.0.0.1:3000`
- Public invoice API: `http://127.0.0.1:6000/invoice?obuid=123`
- Prometheus UI: `http://127.0.0.1:9090`
- Kafka for local development: `127.0.0.1:9092`

### Internal-only endpoints in Compose

- Aggregator HTTP: `http://aggregator:3001/invoice`
- Aggregator gRPC: `aggregator:3002`
- Aggregator metrics: `aggregator:2112`

The aggregator is intentionally not published to the host in Compose, so invoice consumers should go through the gateway.

## Available Commands

```bash
make obu       # run the OBU simulator
make receiver  # run the data receiver
make calc      # run the distance calculator
make agg       # run the aggregator
make gate      # run the gateway
make proto     # regenerate protobuf and gRPC Go files
```

## Public API

Get an invoice through the gateway:

```bash
curl "http://127.0.0.1:6000/invoice?obuid=123"
```

Example response:

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

## Internal APIs

The aggregator's internal HTTP API supports:

- `GET /invoice?obuid=<id>`
- `POST /aggregate`

The protobuf contract for internal gRPC writes lives in `types/ptypes.proto`.

The `distance_calculator` currently uses:

- `Aggreagator.AggregateDistance(DistanceRequest) returns (None)`

## Metrics

For local host runs, the aggregator exposes Prometheus metrics at:

```text
http://127.0.0.1:2112/metrics
```

Prometheus config files:

- `config/prometheus.yml` for local host runs
- `config/prometheus.docker.yml` for Docker Compose

If you have a local Prometheus binary, start it with:

```bash
./prometheus --config.file=./config/prometheus.yml
```

## Current Implementation Notes

These details still matter if you use this repository as a reference or starting point:

- Aggregated distance is stored only in memory. Restarting the aggregator clears all totals.
- The distance calculator keeps a single previous point for the whole process, not one previous point per OBU. Consecutive messages from different vehicles can affect each other's calculated distance.
- The OBU simulator generates random coordinates and random OBU IDs on startup, so values are not deterministic across runs.
- There are no automated test files in the repository today.
- There is still no authentication, persistence layer, or dynamic service discovery. The gateway/internal split is enforced by bind addresses and Docker networking, not by auth.

## Shared Data Types

The main shared payloads are:

- `types.OBUdata`: `{ obuid, lat, long }`
- `types.Distance`: `{ obuid, value, unix }`
- `types.Invoice`: `{ obuid, total_distance, total_amount }`
