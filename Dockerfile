# static-quorum-system
#
# Version 1

# Base on the Fedora image created by Matthew
FROM mattdm/fedora

MAINTAINER Mateus Braga <mateus.a.braga@gmail.com>

# Update packages
#RUN yum update -y

# Install dependencies
RUN yum install -y gcc golang git

#set GOPATH
ENV GOPATH /go

# install static-quorum-system server
RUN go get github.com/mateusbraga/static-quorum-system/...

# By default, launch static-quorum-system server on port 5000
CMD ["/go/bin/static-quorum-systemd", "-bind", ":5000"]

#expose static-quorum-system port
EXPOSE 5000

