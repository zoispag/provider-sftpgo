FROM alpine:3.24.3
RUN apk --no-cache add ca-certificates

ARG TARGETOS
ARG TARGETARCH

ADD "bin/${TARGETOS}_${TARGETARCH}/provider" /usr/local/bin/provider

ENV USER_ID=65532

# The SFTPGo Terraform provider is linked into the binary and runs in-process, so the image ships
# no Terraform CLI or provider plugin. The controller still reads these to record which upstream
# versions it was generated from.
ARG TERRAFORM_VERSION
ARG TERRAFORM_PROVIDER_SOURCE
ARG TERRAFORM_PROVIDER_VERSION
ENV TERRAFORM_VERSION=${TERRAFORM_VERSION}
ENV TERRAFORM_PROVIDER_SOURCE=${TERRAFORM_PROVIDER_SOURCE}
ENV TERRAFORM_PROVIDER_VERSION=${TERRAFORM_PROVIDER_VERSION}

USER ${USER_ID}
EXPOSE 8080

ENTRYPOINT ["provider"]
