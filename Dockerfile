# build stage
FROM golang:1.23.2 AS builder

WORKDIR /app

ADD . .

RUN go mod tidy && cd cmd && go build -o message-service .


# run stage
FROM golang:1.23.2

ENV go_service_env prod

WORKDIR /app

# we need to copy this file to find out the project root
COPY --from=builder /app/go.mod .
COPY --from=builder /app/cmd/message-service .
COPY --from=builder /app/configs/application.yaml ./configs/application.yaml
COPY --from=builder /app/configs/application.local.yaml ./configs/application.local.yaml
COPY --from=builder /app/configs/application.prod.yaml ./configs/application.prod.yaml
COPY --from=builder /app/configs/application.test.yaml ./configs/application.test.yaml

CMD ["./message-service"]
