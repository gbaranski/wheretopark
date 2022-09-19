FROM rust:latest AS builder

RUN apt update && apt install -y libopencv-dev clang libclang-dev libssl-dev

WORKDIR /app
COPY Cargo.toml ./
COPY Cargo.lock ./
COPY providers/cctv ./providers/cctv
COPY rust rust

RUN cargo build --release --bin wheretopark-provider-cctv

FROM debian:buster-slim AS runtime
RUN apt update && apt install -y libopencv-dev clang libclang-dev libssl-dev

WORKDIR /app
COPY --from=builder /app/target/release/wheretopark-provider-cctv /usr/local/bin
ENTRYPOINT ["/usr/local/bin/wheretopark-provider-cctv"]