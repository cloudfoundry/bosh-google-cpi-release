FROM ubuntu:jammy

# Update base image
RUN apt-get update && apt-get -y install locales apt-utils
RUN localedef -i en_US -f UTF-8 en_US.UTF-8

RUN apt-get update -y \
    && apt-get install -y --no-install-recommends \
        build-essential \
        apt-transport-https \
        ca-certificates \
        gnupg \
        git-core \
        curl \
        git \
        wget \
        jq \
        openssh-client \
        ruby \
        sudo \
        sqlite3 \
        libsqlite3-dev \
        mysql-client \
        libmysqlclient-dev \
        libpq-dev \
        libyaml-dev

RUN export GCLOUD_KEYFILE="/usr/share/keyrings/cloud.google.gpg" \
    && echo "deb [signed-by=${GCLOUD_KEYFILE}] https://packages.cloud.google.com/apt cloud-sdk main" \
    | tee -a /etc/apt/sources.list.d/google-cloud-sdk.list \
    && curl -fSL https://packages.cloud.google.com/apt/doc/apt-key.gpg \
    | apt-key --keyring "${GCLOUD_KEYFILE}" add - \
    && apt-get update -y \
    && apt-get install -y --no-install-recommends \
      google-cloud-cli

# Install Golang
COPY --from=bosh/golang-release /var/vcap/packages/golang-1-linux /usr/local/go
ENV GOROOT=/usr/local/go PATH=/usr/local/go/bin:$PATH

# Install Bosh cli
RUN wget https://github.com/cloudfoundry/bosh-cli/releases/download/v7.8.0/bosh-cli-7.8.0-linux-amd64
RUN chmod +x bosh-cli-*
RUN mv bosh-cli-* /usr/bin/bosh

# receipt generator
RUN cd /tmp && \
    curl -o certify-artifacts -L https://s3.amazonaws.com/bosh-certification-generator-releases/certify-artifacts-linux-amd64 && \
    chmod +x certify-artifacts && \
    mv certify-artifacts /bin/certify-artifacts

RUN useradd non-root-user --shell /bin/bash --create-home
RUN echo "non-root-user ALL=(ALL) NOPASSWD: ALL" >> /etc/sudoers
