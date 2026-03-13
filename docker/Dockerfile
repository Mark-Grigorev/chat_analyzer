# defaults
ARG GO_VERSION=1.23.2
ARG ALPINE_VERSION=3.19
ARG VERSION=DEV
ARG PORT=8080

FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS build

ARG VERSION
ARG PATH_TO_MAIN=cmd/*.go

WORKDIR /app

RUN apk add --no-cache gcc libc-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -a -installsuffix cgo -o app \
    -ldflags "-X 'main.version=${VERSION}'" ${PATH_TO_MAIN}

FROM alpine:${ALPINE_VERSION}

ARG PORT

RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup && \
    apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=build --chown=appuser:appgroup /app/app /app/

USER appuser
EXPOSE ${PORT}

CMD ["./app"]
