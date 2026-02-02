ARG BUILDER_IMAGE
FROM ${BUILDER_IMAGE} AS builder

WORKDIR /app
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor \
    -o /bank-account ./services/account/cmd/main.go

FROM alpine:3.18
WORKDIR /root/
COPY --from=builder /bank-account .
COPY .env .
CMD ["./bank-account"]