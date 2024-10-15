#FROM golang:bookworm
#FROM cloudfoundry/cflinuxfs4
FROM jadud/builder

# Copy in the binaries
WORKDIR /app
COPY . .
WORKDIR /app/cmd/queue-server
RUN make container_build
RUN chmod 755 /app/cmd/queue-server/run.sh

CMD ["/app/cmd/queue-server/run.sh"]


