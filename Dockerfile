FROM golang:1.17.3 as dev_build
RUN mkdir /pod_replicator
WORKDIR /pod_replicator
COPY go.mod go.sum /pod_replicator/
COPY . /pod_replicator/
ENV CGO_ENABLED=0
RUN make all
RUN chmod +x pod_replicator

FROM alpine:3.15.0
COPY --from=dev_build /pod_replicator/pod_replicator /usr/bin/pod_replicator
ENTRYPOINT ["pod_replicator"]
CMD ["--help"]
