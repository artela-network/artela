FROM golang:1.21.5-bullseye as build-env

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
COPY ./artela-evm ./artela-evm
COPY ./aspect-runtime ./aspect-runtime
COPY ./aspect-core ./aspect-core

# Reset the working directory for the build
WORKDIR /go/src/github.com/artela-network/artela

# disable optimisation and strip for remote debugging
ENV COSMOS_BUILD_OPTIONS "nostrip,nooptimization"

# build artelad
RUN make build

# Final image
FROM golang:1.21.5-bullseye as final

WORKDIR /

# Add GO Bin to PATH
ENV PATH "/go/bin:${PATH}"

# Copy over binaries from the build-env
COPY --from=build-env /go/src/github.com/artela-network/artela/build/artelad /
COPY --from=build-env /go/src/github.com/artela-network/artela/scripts/start-artela.sh /

# Run artelad by default, omit entrypoint to ease using container with artelad
ENTRYPOINT ["/bin/bash", "-c"]