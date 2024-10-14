#!/bin/bash

# This runs in the container context only.
export VCAP_SERVICES=$(cat /app/vcap.json)

############################################################
# ENV VARS
############################################################
# Should we use service-name, or instance-name?
# The app now creates its own buckets as needed.
# for bucket_name in $(echo "$VCAP_SERVICES" | jq .s3 | jq .[].instance_name);
# do 
#     export eb_bucket="$(echo "$VCAP_SERVICES" | jq .s3 | \
#         jq --raw-output '.[] | select(.instance_name=="database-storage").credentials.bucket')"
#     export eb_uri="$(echo "$VCAP_SERVICES" | jq .s3 | \
#         jq --raw-output '.[] | select(.instance_name=="database-storage").credentials.uri')"
#     export eb_aki="$(echo "$VCAP_SERVICES" | jq .s3 | \
#         jq --raw-output '.[] | select(.instance_name=="database-storage").credentials.access_key_id')"
#     export eb_sak="$(echo "$VCAP_SERVICES" | jq .s3 | \
#         jq --raw-output '.[] | select(.instance_name=="database-storage").credentials.secret_access_key')"

#     ############################################################
#     # CREATE BUCKETS
#     ############################################################
#     # mc [GLOBALFLAGS] alias set \
#     #                  [--api "string"]                           \
#     #                  [--path "string"]                          \
#     #                  ALIAS                                      \
#     #                  URL                                        \
#     #                  ACCESSKEY                                  \
#     #                  SECRETKEY
#     echo Creating bucket $eb_bucket
#     /minio/mc alias set minnie ${eb_uri} ${eb_aki} ${eb_sak}
#     /minio/mc mb minnie/${eb_bucket}
# done

./indexer.exe