FROM --platform=$BUILDPLATFORM golang:1.25-alpine AS build

ARG TARGETOS
ARG TARGETARCH
ARG TARGETVARIANT

WORKDIR /app

COPY . .
RUN go mod download

RUN CGO_ENABLED=0 \
    GOOS=${TARGETOS} \
    GOARCH=${TARGETARCH} \
    GOARM=$(echo ${TARGETVARIANT} | tr -d 'v') \
    go build -o ./bin/huawei-solar-to-influx ./cmd/main.go


FROM alpine:latest AS final
LABEL maintainer="mosmo"

WORKDIR /app

COPY --from=build /app/bin/huawei-solar-to-influx ./

ENTRYPOINT [ "./huawei-solar-to-influx" ]