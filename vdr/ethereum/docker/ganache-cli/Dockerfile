FROM trufflesuite/ganache-cli:latest
RUN /bin/sh -c "apk add --no-cache bash"
ENTRYPOINT node /app/ganache-core.docker.cli.js --quiet --networkId 57771 --verbose --host 0.0.0.0

