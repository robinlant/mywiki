FROM golang:1.25 AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY *.go .

RUN CGO_ENABLED=0 go build -o main .


FROM alpine:latest

WORKDIR /app

COPY *.csv .

COPY --from=build /app/main .

RUN chmod +x main && \
    adduser -HD container && \
    chown -R container /app

USER container

CMD ["./main"]
