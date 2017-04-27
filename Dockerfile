FROM alpine:3.5
LABEL maintainer "Mario Freitas <imkira@gmail.com>"

WORKDIR /usr/local/bin

ENV RELEASE_VER  v0.0.1
ENV RELEASE_SHA1 d46c68802398609213b9ebb9e125e903be099f12

RUN apk add --no-cache --update \
      ca-certificates \
      curl \
    && curl -Lsj -o gcp-iap-auth https://github.com/imkira/gcp-iap-auth/releases/download/${RELEASE_VER}/gcp-iap-auth-linux-amd64 \
    && (echo "${RELEASE_SHA1} *gcp-iap-auth" | sha1sum -c -) \
    && chmod +x gcp-iap-auth \
    && apk del curl \
    && rm -rf /var/cache/apk/*

#ENV GCP_IAP_AUTH_LISTEN_ADDR=0.0.0.0
#ENV GCP_IAP_AUTH_LISTEN_PORT=80
#ENV GCP_IAP_AUTH_AUDIENCES=https://domain1,https://domain2
#ENV GCP_IAP_AUTH_PUBLIC_KEYS=/path/to/public_keys_file
#ENV GCP_IAP_AUTH_TLS_CERT=/path/to/tls_cert_file
#ENV GCP_IAP_AUTH_TLS_KEY=/path/to/tls_key_file

EXPOSE 80 443
ENTRYPOINT ["/usr/local/bin/gcp-iap-auth"]
