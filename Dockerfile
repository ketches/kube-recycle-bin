FROM golang AS builder
ARG KRB_APPNAME

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
COPY cmd/${KRB_APPNAME} ./cmd/${KRB_APPNAME}
COPY pkg/ ./pkg/
COPY internal/ ./internal/

RUN CGO_ENABLED=0 GOOS=linux go build -o ./bin/${KRB_APPNAME} ./cmd/${KRB_APPNAME}

FROM alpine:latest
ARG KRB_APPNAME

WORKDIR /run

COPY --from=builder /app/bin/${KRB_APPNAME} app

ENTRYPOINT ["./app"]