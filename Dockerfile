# https://www.cloudreach.com/en/resources/blog/cts-build-golang-dockerfiles/
# Build
FROM golang:latest as BUILDER

RUN mkdir /build
ADD . /build/
WORKDIR /build

# Build application
RUN make clean; make build

# Run
FROM scratch

COPY --from=BUILDER /build/datanerd /app/
WORKDIR /app

CMD ["./crazy-lemur"]
