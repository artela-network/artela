FROM golang:1.20-bullseye as build-env

# Install minimum necessary dependencies
ENV PACKAGES curl make git libc-dev bash gcc
RUN apt-get update && apt-get upgrade -y && \
    apt-get install -y $PACKAGES

# Set working directory for the source copy
WORKDIR /go/src/github.com/artela-network

# Add source files
COPY ./artela ./artela
COPY ./artela-cosmos-sdk ./artela-cosmos-sdk
COPY ./artela-cometbft ./artela-cometbft

# Reset the working directory for the build
WORKDIR /go/src/github.com/artela-network/artela

# build Ethermint
RUN make build-linux

# Final image
FROM golang:1.20-bullseye as final

WORKDIR /

RUN apt-get update

# Copy over binaries from the build-env
COPY --from=build-env /go/src/github.com/artela-network/artela/build/artelad /
COPY --from=build-env /go/src/github.com/artela-network/artela/scripts/start-artela.sh /

EXPOSE 26656 26657 1317 9090 8545 8546

# Run artelad by default, omit entrypoint to ease using container with artelad
ENTRYPOINT ["/bin/bash", "-c"]