#FROM golang:bookworm
#FROM cloudfoundry/cflinuxfs4
FROM jadud/builder

# Copy in the binaries
WORKDIR /app
COPY . .
WORKDIR /app/cmd/crawler
RUN make container_build
RUN chmod 755 /app/cmd/crawler/run.sh

CMD ["/app/cmd/crawler/run.sh"]


