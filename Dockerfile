FROM registry.access.redhat.com/ubi9/go-toolset:1.18.10-4 as build
# args required for image labels
ARG BASE_IMAGE_DIGEST
ARG BASE_IMAGE_NAME
ARG BUILD_DATE
ARG COMMIT_HASH
ARG CRDA_VERSION
USER root
# build crda cli
WORKDIR /crda
COPY . .
# build and sanitize file name (ie crda-1.2.3-linux-amd64 >>> crda)
RUN make build \
    && mv "$(find build -name 'crda*')" build/crda
USER default


 # if the base image here is modified, modify BASE_IMAGE_NAME in Makefile as well
FROM registry.access.redhat.com/ubi9/go-toolset:1.18.10-4
# the go-toolset:1.18.10-4 image comes with:
# - go 1.18.10
# - python 3.9.14
# - nodejs 16.18.1 (npm 8.19.2)
# - java and maven are installed manually later

# args required for image labels
ARG BASE_IMAGE_DIGEST
ARG BASE_IMAGE_NAME
ARG BUILD_DATE
ARG COMMIT_HASH
ARG CRDA_VERSION
USER root
# install jdk 17
RUN yum update -y  \
    && yum install java-17-openjdk -y \
    && yum clean all
# install maven 3.9.1 \
RUN curl https://dlcdn.apache.org/maven/maven-3/3.9.1/binaries/apache-maven-3.9.1-bin.tar.gz --output apache-maven-3.9.1-bin.tar.gz \
    && tar xzvf apache-maven-3.9.1-bin.tar.gz \
    && mv apache-maven-3.9.1 /usr/share/maven \
    && ln -s /usr/share/maven/bin/mvn /usr/bin/mvn \
    && rm -r apache-maven-3.9.1*
# symlinks for python3 and pip3 >> python and pip
RUN ln -s /usr/bin/python3 /usr/bin/python && \
    ln -s /usr/bin/pip3 /usr/bin/pip
# copy crda cli
COPY --from=build /crda/build/crda /usr/bin/crda
COPY --from=build /crda/LICENSE /licenses/crda-license

ENTRYPOINT ["/usr/bin/crda", "--client", "image", "--no-color", "--json=true"]

WORKDIR /app
USER default

LABEL org.opencontainers.image.created=$BUILD_DATE \
org.opencontainers.image.authors="Ecosystem Engineering Team, Red Hat" \
org.opencontainers.image.url="https://quay.io/repository/ecosystem-appeng/crda-cli" \
org.opencontainers.image.documentation="https://github.com/RHEcosystemAppEng/crda-cli" \
org.opencontainers.image.source="https://github.com/RHEcosystemAppEng/crda-cli.git" \
org.opencontainers.image.version=$CRDA_VERSION \
org.opencontainers.image.revision=$COMMIT_HASH \
org.opencontainers.image.vendor="Red Hat, Inc." \
org.opencontainers.image.licenses="Apache-2.0" \
org.opencontainers.image.ref.name=$CRDA_VERSION \
org.opencontainers.image.title="Crda CLI" \
org.opencontainers.image.description="Create CodeReady Dependency Analytics Reports" \
org.opencontainers.image.base.digest=$BASE_IMAGE_DIGEST \
org.opencontainers.image.base.name=$BASE_IMAGE_NAME
