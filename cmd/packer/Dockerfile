#FROM golang:bookworm
#FROM cloudfoundry/cflinuxfs4
FROM jadud/builder

# Copy in the binaries
WORKDIR /app
COPY . .
WORKDIR /app/cmd/packer
RUN GOBIN=/bin go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# RUN wget https://downloads.sqlc.dev/sqlc_1.27.0_linux_amd64.tar.gz \ 
#     && tar xvzf sqlc_1.27.0_linux_amd64.tar.gz \
#     && chmod 755 sqlc \
#     && mv sqlc /bin/sqlc

RUN make container_build
RUN cp /app/pkg/dyndb/schema.sql /app/cmd/packer/schema.sql
RUN chmod 755 /app/cmd/packer/run.sh

CMD ["/app/cmd/packer/run.sh"]


