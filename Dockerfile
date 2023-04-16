FROM registry.access.redhat.com/ubi9/ubi-minimal:9.1 AS builder
ARG TARGETARCH
USER root
RUN microdnf install -y tar gzip make which

# install platform specific go version
RUN curl -O -J  https://dl.google.com/go/go1.19.6.linux-${TARGETARCH}.tar.gz
RUN tar -C /usr/local -xzf go1.19.6.linux-${TARGETARCH}.tar.gz
RUN ln -s /usr/local/go/bin/go /usr/local/bin/go

WORKDIR /workspace

COPY . ./
RUN CGO_ENABLED=0 GO111MODULE=on

RUN go mod tidy
RUN go build -o arch-go .

FROM registry.access.redhat.com/ubi9/ubi-minimal:9.1
COPY --from=builder /usr/local/go/ /usr/local/go/
RUN ln -s /usr/local/go/bin/go /usr/local/bin/go
RUN CGO_ENABLED=0 GO111MODULE=on
RUN mkdir /.cache && chmod 777 /.cache
RUN mkdir /go && chmod 777 /go

COPY --from=builder /workspace/arch-go /usr/local/bin/
USER 1001
ENTRYPOINT ["/usr/local/bin/arch-go"]

LABEL name="arch-go" \
      version="0.0.1" \
      summary="arch-go service" \
      description="arch-go checks if your software architecture complies to a given ruleset"
