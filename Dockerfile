FROM alpine:3.5
LABEL maintainer "Mario Freitas <imkira@gmail.com>"

WORKDIR /usr/local/bin

COPY dist/gcp-iap-auth-linux-amd64 /usr/local/bin/gcp-iap-auth

#ENV GCP_IAP_AUTH_LISTEN_ADDR=0.0.0.0
#ENV GCP_IAP_AUTH_LISTEN_PORT=80
#ENV GCP_IAP_AUTH_AUDIENCES=https://domain1,https://domain2
#ENV GCP_IAP_AUTH_PUBLIC_KEYS=/path/to/public_keys_file
#ENV GCP_IAP_AUTH_TLS_CERT=/path/to/tls_cert_file
#ENV GCP_IAP_AUTH_TLS_KEY=/path/to/tls_key_file

EXPOSE 80 443
ENTRYPOINT ["/usr/local/bin/gcp-iap-auth"]
