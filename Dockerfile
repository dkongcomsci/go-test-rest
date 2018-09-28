FROM golang:1.8

ADD /instantclient_12_1 /usr/instantclient_12_1

RUN apt-get update
RUN apt-get install -y pkg-config
RUN apt-get install -y libaio1

ADD oci8.pc /usr/lib/pkgconfig/oci8.pc

ENV LD_LIBRARY_PATH /usr/lib:/usr/local/lib:/usr/instantclient_12_1

RUN ln -s /usr/instantclient_12_1/libclntsh.so.12.1 /usr/instantclient_12_1/libclntsh.so
RUN ln -s /usr/instantclient_12_1/libclntshcore.so.12.1 /usr/instantclient_12_1/libclntshcore.so
RUN ln -s /usr/instantclient_12_1/libocci.so.12.1 /usr/instantclient_12_1/libocci.so

USER nobody

RUN mkdir -p /go/src/github.com/dkongcomsci/go-test-rest
WORKDIR /go/src/github.com/dkongcomsci/go-test-rest

COPY . /go/src/github.com/dkongcomsci/go-test-rest
RUN go-wrapper download && go-wrapper install

CMD ["go-wrapper", "run"]
