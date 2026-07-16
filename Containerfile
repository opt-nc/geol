FROM debian:13.6-slim
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates && rm -rf /var/lib/apt/lists/*
ARG TARGETPLATFORM
COPY ${TARGETPLATFORM}/geol /usr/local/bin/geol
ENV CLICOLOR_FORCE=1
ENV COLORTERM=truecolor
ENTRYPOINT ["geol"]
