FROM cloudfoundry/cflinuxfs4

RUN apt-get update \
    && apt-get install -y \
        build-essential \
        curl \
        software-properties-common \
        golang-go \
        tree \
        wget 

RUN add-apt-repository ppa:longsleep/golang-backports
RUN apt-get update \
    && apt-get install -y \
        golang-go

WORKDIR /minio
RUN curl https://dl.min.io/client/mc/release/linux-amd64/mc \
        --create-dirs \
        -o /minio/mc
RUN chmod +x /minio/mc

WORKDIR /sqlc
RUN wget https://downloads.sqlc.dev/sqlc_1.27.0_linux_amd64.tar.gz \ 
    && tar xvzf sqlc_1.27.0_linux_amd64.tar.gz \
    && chmod 755 sqlc \
    && mv sqlc /bin/sqlc


ENTRYPOINT ["tree /"]