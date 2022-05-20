FROM golang:latest as base

RUN useradd -m praveen


FROM base as dev

RUN apt-get update 

# image processing library
RUN apt-get install -y libvips-dev --fix-missing     


WORKDIR /src/app
COPY . .

RUN chown praveen /src/app

USER praveen




EXPOSE 4200

EXPOSE 25 465 587 2525


RUN go build -tags netgo -ldflags '-s -w' -o ./build/app cmd/main.go

CMD [ "./build/app" ]