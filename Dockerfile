FROM golang:1.22-bookworm AS builder

ARG SERVICE

RUN apt-get update && apt-get install -y --no-install-recommends \
	build-essential \
	ca-certificates \
	librdkafka-dev \
	pkg-config \
	&& rm -rf /var/lib/apt/lists/*

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN test -n "$SERVICE"
RUN CGO_ENABLED=1 GOOS=linux go build -o /out/service ./${SERVICE}


FROM debian:bookworm-slim AS runtime

RUN apt-get update && apt-get install -y --no-install-recommends \
	ca-certificates \
	librdkafka1 \
	&& rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /out/service /app/service

ENTRYPOINT ["/app/service"]
