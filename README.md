# gcp-iap-auth

[![License](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://github.com/imkira/gcp-iap-auth/blob/master/LICENSE.txt)
[![Build Status](http://img.shields.io/travis/imkira/gcp-iap-auth.svg?style=flat)](https://travis-ci.org/imkira/gcp-iap-auth)

`gcp-iap-auth` is a simple server implementation and package in
[Go](http://golang.org) for helping you secure your web apps running on GCP
behind a
[Google Cloud Platform's IAP (Identity-Aware Proxy)](https://cloud.google.com/iap/docs/) by validating IAP signed headers in the requests.

# Why

Validating signed headers helps you protect your app from the following kinds of risks:

- IAP is accidentally disabled;
- Misconfigured firewalls;
- Access from within the project.

## How to use it as a package

```go
go get -u github.com/imkira/gcp-iap-auth/jwt
```

The following is just an excerpt of the provided [simple.go example](https://github.com/imkira/gcp-iap-auth/tree/master/examples):

```go
// Here we validate the tokens in all requests going to
// our server at http://127.0.0.1:12345/auth
// For valid tokens we return 200, otherwise 401.
func AuthHandler(w http.ResponseWriter, req *http.Request) {
	if err := jwt.ValidateRequestClaims(req, cfg); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}
```

For advanced usage, make sure to check the
[available documentation here](http://godoc.org/github.com/imkira/gcp-iap-auth).

## How to use it as a server

[Binary Releases](https://github.com/imkira/releases) are provided for convenience.

After downloading it, you can execute it like:

```shell
gcp-iap-auth --audiences=https://APP_DOMAIN
```

HTTPS is also supported. Just make sure you give it the cert/key files:

```shell
gcp-iap-auth --audiences=https://APP_DOMAIN --tls-cert=PATH_TO_CERT_FILE --tls-key=PATH_TO_KEY_FILE
```

It is also possible to use environment variables instead of flags.
Just prepend `GCP_IAP_AUTH_` to the flag name (in CAPS and with `-` replaced by `_`) and you're good to go (eg: `GCP_IAP_AUTH_AUDIENCES` replaces `--audiences`)

For help, just check usage:

```shell
gcp-iap-auth --help
```

## Integration with NGINX

You can also integrate `gcp-iap-auth` server with [NGINX](https://nginx.org)
using the
[http_auth_request_module](https://nginx.org/en/docs/http/ngx_http_auth_request_module.html).

The important part is as follows ([full nginx.conf example file here](https://github.com/imkira/gcp-iap-auth/tree/master/examples)):

```
    upstream AUTH_SERVER_UPSTREAM {
      server AUTH_SERVER_ADDR:AUTH_SERVER_PORT;
    }

    upstream APP_SERVER_UPSTREAM {
      server APP_SERVER_ADDR:APP_SERVER_PORT;
    }

    server {
      server_name APP_DOMAIN;

      location = /gcp-iap-auth {
          internal;
          proxy_pass                 http://AUTH_SERVER_UPSTREAM/auth;
          proxy_pass_request_body    off;
          proxy_pass_request_headers off;
          proxy_set_header           X-Goog-Authenticated-User-JWT $http_x_goog_authenticated_user_jwt;
      }

      location / {
        auth_request /gcp-iap-auth;
        proxy_pass   http://APP_SERVER_UPSTREAM;
      }
    }
```

Please note:

- Replace `AUTH_SERVER_UPSTREAM`, `AUTH_SERVER_ADDR`, and `AUTH_SERVER_PORT` with the data about your `gcp-iap-auth` server.
- Replace `APP_SERVER_UPSTREAM`, `APP_SERVER_ADDR`, and `APP_SERVER_PORT` with the data about your own web app server.
- Replace `APP_DOMAIN` with the domain(s) you set up in your GCP IAP settings.
- `gcp-iap-auth` only needs to receive the original `X-Goog-Authenticated-User-JWT` header sent by Google, so you can and you are advised to disable proxying the original request body and other headers. Not only it is unecessary you may leak information you may not want to.
- Please adjust appropriately (you may want to use HTTPS instead of HTTP, multiple domains, etc.). This example is just provided for reference.

## License

gcp-iap-auth is licensed under the MIT license:

www.opensource.org/licenses/MIT

## Copyright

Copyright (c) 2017 Mario Freitas. See
[LICENSE](https://github.com/imkira/gcp-iap-auth/blob/master/LICENSE.txt)
for further details.
