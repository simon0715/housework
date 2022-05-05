FROM centos
WORKDIR /usr/local
ADD go1.18.1.linux-amd64.tar.gz .
ENV GOROOT /usr/local/go
ENV PATH $PATH:$GOROOT/bin
RUN mkdir gowork
WORKDIR gowork
ADD housework-master/main.go .
ENTRYPOINT ["go","run","main.go"]

