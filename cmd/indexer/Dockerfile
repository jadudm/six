#FROM golang:bookworm
#FROM cloudfoundry/cflinuxfs4
FROM jadud/builder

# Copy in the binaries
# FIXME: /home/vcap/app to mirror deploy
WORKDIR /app
COPY . .
WORKDIR /app/cmd/scraper
RUN make container_build
RUN chmod 755 /app/cmd/scraper/run.sh

CMD ["/app/cmd/scraper/run.sh"]


