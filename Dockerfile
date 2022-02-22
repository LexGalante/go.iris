# see https://hub.docker.com/_/golang/
FROM golang:1.17-alpine
# define homework
WORKDIR /app
# define default environments variables
ENV APP_SECRET=LEXGALANTE.GO.IRIS \
    APP_PORT=5050 \
    DB_NAME=tracking \    
    DB_URI=mongodb://root:mongo@localhost:27017/
# send module dependencies to container
COPY go.mod ./
COPY go.sum ./
# go get dependencies
RUN go mod download
# copy all file .go to container
COPY *.go ./
# prepare compiled program
RUN go build -o /program
# define entrtpoint of application
CMD [ "/program" ]

