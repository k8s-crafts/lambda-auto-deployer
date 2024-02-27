FROM docker.io/library/golang:1.20 as build

WORKDIR /build
# Copy dependencies list

COPY go.mod go.sum ./

# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

COPY main.go .
COPY utils ./utils

# Build with optional lambda.norpc tag
# to remove the additional process and RPC hop, resulting in faster invokes
RUN GOOS=linux GOARCH=amd64 go build -tags lambda.norpc -o app main.go

FROM public.ecr.aws/lambda/provided:al2023

COPY --from=build /build/app ./app

ENTRYPOINT [ "./app" ]
