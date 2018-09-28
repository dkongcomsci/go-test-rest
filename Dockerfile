# Download the -basic, -basiclite and -devel RPMs from http://www.oracle.com/technetwork/topics/linuxx86-64soft-092277.html
# and put them in this dir.
#
# Make sure oci8.pc matches the InstantClient Version
# Optionally setup tnsnames.ora and sqlnet.ora. You can still use ezconnect syntax without those files

#
# Create build image
#
FROM golang:latest

# Install Oracle Instantclient
WORKDIR /tmp
COPY *.rpm ./
COPY ./oci8.pc /usr/share/pkgconfig/
RUN yum -y install oracle-instantclient12.2-basic-12.2.0.1.0-1.x86_64.rpm oracle-instantclient12.2-devel-12.2.0.1.0-1.x86_64.rpm
ENV ORACLE_HOME=/usr/lib/oracle/12.2/client64
ENV PATH=$PATH:$ORACLE_HOME/bin
ENV LD_LIBRARY_PATH=$ORACLE_HOME/lib
#
# Build the example app. Change this bit to suit
#
RUN mkdir -p /root/go/src/myapp
WORKDIR /root/go/src/myapp
COPY *.go .
# RUN go get github.com/mattn/go-oci8
RUN go get gopkg.in/rana/ora.v4
RUN go build -ldflags="-s -w" .

#
# Create run image
#
FROM centos:latest
LABEL maintainer="Dalibor Andzakovic <dali@swerve.nz>" instantclient="12.2"

# Update image
RUN yum -y update
# Install Oracle Instantclient
WORKDIR /tmp
COPY *.rpm ./
RUN yum -y install install oracle-instantclient12.2-basiclite-12.2.0.1.0-1.x86_64.rpm
RUN rm /tmp/*.rpm && yum clean all
ENV ORACLE_HOME=/usr/lib/oracle/12.2/client64
ENV PATH=$PATH:$ORACLE_HOME/bin
ENV LD_LIBRARY_PATH=$ORACLE_HOME/lib
RUN mkdir -p ${ORACLE_HOME}/network/admin
#COPY tnsnames.ora ${ORACLE_HOME}/network/admin/
#COPY sqlnet.ora ${ORACLE_HOME}/network/admin/

# Copy myapp
USER nobody

RUN mkdir -p /go/src/github.com/dkongcomsci/go-test-rest
WORKDIR /go/src/github.com/dkongcomsci/go-test-rest

COPY . /go/src/github.com/dkongcomsci/go-test-rest
RUN go-wrapper download && go-wrapper install

CMD ["go-wrapper", "run"]