ARG APP_NAME=reckon

FROM rust:1.51.0-buster as builder

ARG APP_NAME
WORKDIR /app/${APP_NAME}

RUN apt-get update \
 && apt-get install -y --no-install-recommends \
    pkg-config \
    libssl-dev \
 && rm -rf /var/lib/apt/lists/*

COPY Cargo.toml Cargo.lock ./
RUN mkdir src \
 && echo 'fn main() {println!("if you see this, the build broke")}' > src/main.rs \
 && cargo build --release \
 && rm -f target/release/deps/${APP_NAME}*

COPY . .
RUN cargo build --release

FROM debian:buster-slim

ARG APP_NAME
ENV APP_NAME=${APP_NAME}
WORKDIR /usr/local/bin

RUN apt-get update \
 && apt-get install -y --no-install-recommends \
    libssl-dev \
    ca-certificates \
 && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/${APP_NAME}/target/release/${APP_NAME} ${APP_NAME}

CMD ${APP_NAME}
