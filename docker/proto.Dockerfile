# STEP 1: Define base images for reuse
FROM bufbuild/buf:1.50.0 AS buf-base
FROM golang:1.25.6-alpine AS builder

# STEP 2: Install dependencies and compile binaries
RUN apk add --no-cache git

# Install Protobuf and gRPC plugins for Go
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest && \
    go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest

# STEP 3: Fetch external proto dependencies (googleapis)
FROM buf-base AS deps-fetcher
WORKDIR /tmp-deps
COPY buf.yaml .
RUN mkdir proto && touch proto/temp.proto

# Export googleapis directly into a local folder to include in the final image
RUN buf dep update && buf export buf.build/googleapis/googleapis --output /proto-includes

# STEP 4: Create the final lightweight image
FROM buf-base

# Copy compiled binaries from the builder stage
COPY --from=builder /go/bin/protoc-gen-go /usr/local/bin/
COPY --from=builder /go/bin/protoc-gen-go-grpc /usr/local/bin/
COPY --from=builder /go/bin/protoc-gen-grpc-gateway /usr/local/bin/

# Copy baked-in dependencies to the internal include path
# This allows 'buf' to find google/api/annotations.proto without internet access
COPY --from=deps-fetcher /proto-includes /usr/local/include

# Set the working directory inside the container
WORKDIR /workspace

# Set 'buf' as the default entrypoint for the container
ENTRYPOINT ["buf"]