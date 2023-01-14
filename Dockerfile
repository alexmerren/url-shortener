# --- Build Image ---

FROM golang:1.18 as build-image

WORKDIR /build

ADD go.mod go.sum Makefile ./
ADD ./internal ./internal
ADD ./cmd ./cmd
ADD ./vendor ./vendor

RUN CGO_ENABLED=0 make build

# --- REST Server Image ---

FROM golang:1.18-alpine as rest-server-image

EXPOSE 8080

WORKDIR /usr/local/

COPY --from=build-image /build/dist/url-shortener /usr/local/dist/url-shortener
