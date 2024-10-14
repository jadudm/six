FROM cloudfoundry/cflinuxfs4

RUN apt-get update \
    && apt-get install -y \
        build-essential \
        curl \
        software-properties-common \
        golang-go \
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
