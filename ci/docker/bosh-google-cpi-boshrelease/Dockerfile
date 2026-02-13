FROM ubuntu

ENV DEBIAN_FRONTEND=noninteractive
ENV TZ=Etc/UTC

ENV BOSH_CLI_VERSION=7.8.0

# Update base image
RUN apt-get update -y \
    && apt-get -y install \
      apt-transport-https \
      apt-utils \
      curl \
      gnupg \
      locales \
    && localedef -i en_US -f UTF-8 en_US.UTF-8 \
    && export GCLOUD_KEYFILE="/usr/share/keyrings/cloud.google.gpg" \
    && echo "deb [signed-by=${GCLOUD_KEYFILE}] https://packages.cloud.google.com/apt cloud-sdk main" \
     | tee -a /etc/apt/sources.list.d/google-cloud-sdk.list \
    && curl -fsSL https://packages.cloud.google.com/apt/doc/apt-key.gpg \
     | apt-key --keyring "${GCLOUD_KEYFILE}" add - \
    && apt-get update -y \
    && apt-get install -y --no-install-recommends \
        build-essential \
        ca-certificates \
        git-core \
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
        libyaml-dev \
        google-cloud-cli

# Install Golang
COPY --from=ghcr.io/cloudfoundry/bosh/golang-release:latest /var/vcap/packages/golang-1-linux /usr/local/go
ENV GOROOT=/usr/local/go PATH=/usr/local/go/bin:$PATH

# Install Bosh cli
RUN bosh_bin_path="/usr/bin/bosh" \
    && curl -fsSL https://github.com/cloudfoundry/bosh-cli/releases/download/v${BOSH_CLI_VERSION}/bosh-cli-${BOSH_CLI_VERSION}-linux-amd64 \
      > "${bosh_bin_path}" \
    && chmod +x "${bosh_bin_path}"

RUN useradd non-root-user --shell /bin/bash --create-home \
    && echo "non-root-user ALL=(ALL) NOPASSWD: ALL" \
      >> /etc/sudoers
