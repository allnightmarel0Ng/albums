FROM golang:1.23.1-alpine3.19 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN apk --no-cache update && \
    apk --no-cache add gcc libc-dev

ENV CGO_ENABLED 1

ARG SERVICE_NAME
RUN go build -tags musl -o main ./cmd/${SERVICE_NAME}

FROM alpine:3.19

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/.env .

ARG SERVICE_PORT
ARG EXPOSE_PORT

RUN if [ -n "$EXPOSE_PORT" ]; then \
        echo "EXPOSE $SERVICE_PORT" >> /app/Dockerfile.tmp; \
    fi

CMD ["./main"]