FROM golang:1.18 AS build

RUN useradd -u 10001 benthos

WORKDIR /build/
COPY . /build/

RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor

FROM busybox AS package

LABEL maintainer="Ashley Jeffs <ash@jeffail.uk>"

WORKDIR /

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /build/benthos-plugin-example .
COPY ./config/example_1.yaml /benthos.yaml

USER benthos

EXPOSE 4195

ENTRYPOINT ["/benthos-plugin-example"]

CMD ["-c", "/benthos.yaml"]
