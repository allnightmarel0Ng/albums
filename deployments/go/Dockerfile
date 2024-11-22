FROM golang:1.23.1-alpine3.19

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

ARG SERVICE_NAME
RUN go build -o main ./cmd/${SERVICE_NAME}

ARG SERVICE_PORT
EXPOSE ${SERVICE_PORT}

CMD ["./main"]