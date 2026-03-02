FROM golang:1.25-alpine AS build

WORKDIR /app

COPY . .
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/practice ./cmd/main.go


FROM alpine:latest AS final
LABEL maintainer="mosmo"

WORKDIR /app

COPY --from=build /app/bin/practice ./

EXPOSE 8080

ENTRYPOINT [ "./practice" ]