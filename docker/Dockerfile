# defaults
ARG GO_VERSION=1.20.1
ARG ALPINE_VERSION=3.19
ARG VERSION=DEV
ARG PORT=8080
ARG PATH_TO_MAIN=cmd/*.go

FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS build
ARG VERSION

WORKDIR /app

RUN  apk add --no-cache gcc libc-dev
COPY go.mod go.sum ./

RUN go mod download && go mod vendor

COPY . .

RUN go build -mod=vendor -a -installsuffix cgo -o app \
    -lgflags "-X 'main.version=${VERSION}}'" ${PATH_TO_MAIN}

FROM alpine:${ALPINE_VERSION}

ARG PORT

RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup && \
    apk add --no-cache ca-certificated tzdata

WORKDIR /app

COPY --from=build --chown=appuser:appgroup /app/app /app/

USER appuser
EXPOSE ${PORT}

CMD ["./app"]