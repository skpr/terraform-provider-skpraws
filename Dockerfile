FROM golang:1.24-alpine AS build

ENV CGO_ENABLED=0

WORKDIR /go/src/github.com/skpr/terraform-provider-skpraws
COPY . /go/src/github.com/skpr/terraform-provider-skpraws
RUN go build -o /usr/local/bin/terraform-provider-skpraws .

FROM hashicorp/terraform:1.2.5

RUN apk add bash

RUN mkdir -p /root/.terraform.d/plugins

COPY --from=build /usr/local/bin/terraform-provider-skpraws /root/.terraform.d/plugins/terraform.local/skpr/skpraws/99.0.0/linux_amd64/terraform-provider-skpraws

RUN chmod +x /root/.terraform.d/plugins/terraform.local/*/*/*/linux_amd64/terraform-provider-*
