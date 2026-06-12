FROM debian:13.5-slim
ARG TARGETPLATFORM
COPY ${TARGETPLATFORM}/geol /usr/local/bin/geol
ENTRYPOINT ["geol"]
