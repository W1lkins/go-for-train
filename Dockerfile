FROM golang:alpine as builder

ENV PATH /go/bin:/usr/local/go/bin:$PATH
ENV GOPATH /go

RUN	apk add --no-cache ca-certificates bash

COPY . /go/src/github.com/evalexpr/go-for-train

RUN set -x \
	&& apk add --no-cache --virtual .build-deps \
		git \
		gcc \
		libc-dev \
		libgcc \
		make \
	&& cd /go/src/github.com/evalexpr/go-for-train \
    && make vendor \
    && make static \
	&& mv go-for-train /usr/bin/go-for-train \
	&& apk del .build-deps \
	&& rm -rf /go \
	&& echo "Build complete."

FROM alpine:latest

COPY --from=builder /usr/bin/go-for-train /usr/bin/go-for-train
COPY --from=builder /etc/ssl/certs/ /etc/ssl/certs

RUN adduser -D -u 1000 user \
  && chown -R user /home/user

RUN mkdir -p /home/user/.config/go-for-train/config

USER user

ENV USER user

WORKDIR /home/user

ENTRYPOINT [ "go-for-train" ]
CMD [ "--help" ]
