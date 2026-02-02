ARG BUILDER_IMAGE
FROM ${BUILDER_IMAGE} AS builder

WORKDIR /app

RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor \
    -o /bank-account ./services/account/cmd/main.go

FROM alpine:3.18
WORKDIR /root/
COPY --from=builder /bank-account .
CMD ["./bank-account"]