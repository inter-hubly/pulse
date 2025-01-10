FROM golang:latest as builder

# Install dependencies
RUN apt-get update && apt-get install -y git openssh-client

# Ensure the SSH directory exists
RUN mkdir -p /root/.ssh

# Check if SSH_KEY is set (for debugging)
RUN echo "SSH_KEY is: $SSH_KEY"

# Decode SSH key and set permissions
ARG SSH_KEY
RUN echo "$SSH_KEY" | base64 --decode > /root/.ssh/id_ed25519 && chmod 600 /root/.ssh/id_ed25519

# Alternatively, fallback if the previous approach fails (debugging)
# RUN mkdir -p /root/.ssh && echo "$SSH_KEY" > /root/.ssh/id_ed25519 && chmod 600 /root/.ssh/id_ed25519
RUN ssh-keyscan github.com >> /root/.ssh/known_hosts
RUN git config --global url.ssh://git@github.com/.insteadOf https://github.com/

ENV GOPRIVATE=github.com/inter-hubly/*

ENV ENVIRONMENT=${ENVIRONMENT}
ENV ENVIRONMENT_PORT=${ENVIRONMENT_PORT}
ENV ENVIRONMENT_HASH_ENCRYPT=${ENVIRONMENT_HASH_ENCRYPT}
ENV ENVIRONMENT_APPID=${ENVIRONMENT_APPID}
ENV ENVIRONMENT_ELASTICSEARCH_HOST=${ENVIRONMENT_ELASTICSEARCH_HOST}
ENV ENVIRONMENT_REDIS_HOST=${ENVIRONMENT_REDIS_HOST}
ENV ENVIRONMENT_REDIS_DATABASE=${ENVIRONMENT_REDIS_DATABASE}
ENV ENVIRONMENT_AMPQ_HOST=${ENVIRONMENT_AMPQ_HOST}
ENV ENVIRONMENT_PGSQL_HOST=${ENVIRONMENT_PGSQL_HOST}

COPY . /app
WORKDIR /app
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o pulse /app/internal/src/main.go

FROM alpine
ARG ENVIRONMENT
LABEL maintainer="pulse"
WORKDIR /app
COPY --from=builder /app/pulse /app
ENTRYPOINT ./pulse
