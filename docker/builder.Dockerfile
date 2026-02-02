FROM golang:1.25.6-alpine

# 1. Install system dependencies
RUN apk add --no-cache git make build-base

WORKDIR /app

# 2. Copy dependency definitions first (Layer 1)
COPY go.mod go.sum* ./


# 3. Download dependencies independently (Cacheable layer)
# This prevents re-downloading unless go.mod or go.sum changes
#RUN go mod download

# 4. Copy the entire source code (Layer 2)
COPY . .

# 5. Run tidy and vendor only if the vendor directory is missing
# RUN if [ ! -d vendor ]; then \
#   go mod tidy && go mod vendor ; \
#   fi
RUN go mod verify

# Keep the container running for development use
CMD ["sleep", "infinity"]