FROM golang:alpine3.21

RUN mkdir /build
WORKDIR /build
#RUN cd /build && git clone https://github.com/saviobarr/GOWEBAPI.git
COPY . /build/
RUN cd /build && go build
EXPOSE 8080
ENTRYPOINT [ "/build/go-api" ]
