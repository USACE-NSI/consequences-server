FROM ghcr.io/osgeo/gdal:ubuntu-full-3.8.3 AS builder

ENV TZ=America/New_York
ENV PATH=/go/bin:$PATH
ENV GOROOT=/go
ENV GOPATH=/src/go

RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone &&\
    mkdir /go &&\
    mkdir -p /src/go &&\
    apt update &&\
    apt -y install build-essential &&\
	apt -y install pkg-config &&\
    apt -y install gdal-bin gdal-data libgdal-dev &&\
    apt -y install wget &&\
    wget https://go.dev/dl/go1.25.8.linux-amd64.tar.gz -P / &&\
    tar -xvzf /go1.25.8.linux-amd64.tar.gz -C / &&\
    apt -y install git


WORKDIR /app

# Fetch the project
RUN git clone https://github.com/USACE-NSI/consequences-server .

# Build the binary
# Note: Using server.go as specified in your latest snippet
RUN go mod download
RUN CGO_ENABLED=1 GOOS=linux go build -o /consequences-server server.go

# Stage 2: Production
FROM ghcr.io/osgeo/gdal:ubuntu-full-3.8.3 AS prod
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone &&\
    apt update &&\
	apt -y install pkg-config &&\
    apt -y install gdal-bin gdal-data libgdal-dev
# Copy the compiled binary from the builder stage
COPY --from=builder /consequences-server /usr/local/bin/consequences-server

EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/consequences-server"]