###First Stage###

#Finds image from AWS public repo.
FROM public.ecr.aws/docker/library/golang:latest AS build
WORKDIR /build
COPY . .
RUN go mod download

#Build go binary for linux.
RUN CGO_ENABLED=0 GOOS=linux go build -a -o main .


###Second stage###
#From AWS public repo.
FROM public.ecr.aws/docker/library/alpine:latest

RUN apk update && \
    apk upgrade && \
    apk add ca-certificates && \
    apk add tzdata

WORKDIR /app

#Copy files from first stage of build
COPY --from=build /build/main ./

RUN pwd && find .

#Identify port to listen on
EXPOSE 8080

#Execute binary.
CMD ["./main"]