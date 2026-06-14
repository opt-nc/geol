FROM debian:13.5-slim
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates && rm -rf /var/lib/apt/lists/*
ARG TARGETPLATFORM
COPY ${TARGETPLATFORM}/geol /usr/local/bin/geol
ENTRYPOINT ["geol"]
