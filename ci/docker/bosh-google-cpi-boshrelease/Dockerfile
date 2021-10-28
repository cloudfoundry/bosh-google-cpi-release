FROM ubuntu:16.04

# Update base image
RUN DEBIAN_FRONTEND=noninteractive apt-get update && apt-get -y install locales apt-utils
RUN localedef -i en_US -f UTF-8 en_US.UTF-8
RUN DEBIAN_FRONTEND=noninteractive apt-get -y upgrade; apt-get clean

# Install dependencies
RUN DEBIAN_FRONTEND=noninteractive apt-get install -y build-essential git-core curl wget jq sudo; apt-get clean
RUN DEBIAN_FRONTEND=noninteractive apt-get install -y sqlite3 libsqlite3-dev; apt-get clean
RUN DEBIAN_FRONTEND=noninteractive apt-get install -y mysql-client libmysqlclient-dev; apt-get clean
RUN DEBIAN_FRONTEND=noninteractive apt-get install -y libpq-dev; apt-get clean
RUN DEBIAN_FRONTEND=noninteractive apt-get install -y python python-pip libyaml-dev libpython-dev; apt-get clean

# Install Golang
COPY --from=bosh/golang-release /var/vcap/packages/golang-1-linux /usr/local/go
ENV GOROOT=/usr/local/go PATH=/usr/local/go/bin:$PATH

# Install Google Cloud SDK
ENV GCLOUD_SDK_VERSION 305.0.0
RUN curl -sSL https://dl.google.com/dl/cloudsdk/channels/rapid/downloads/google-cloud-sdk-${GCLOUD_SDK_VERSION}-linux-x86_64.tar.gz | tar -v -C /usr/local -xz
ENV PATH $PATH:/usr/local/google-cloud-sdk/bin

# Instal chruby
RUN mkdir /tmp/chruby && \
    cd /tmp && \
    curl https://codeload.github.com/postmodern/chruby/tar.gz/v0.3.9 | tar -xz && \
    cd /tmp/chruby-0.3.9 && \
    sudo ./scripts/setup.sh && \
    rm -rf /tmp/chruby

# Install ruby-install
RUN mkdir /tmp/ruby-install && \
    cd /tmp && \
    curl https://codeload.github.com/postmodern/ruby-install/tar.gz/v0.7.0 | tar -xz && \
    cd /tmp/ruby-install-0.7.0 && \
    sudo make install && \
    rm -rf /tmp/ruby-install


# Set default ruby
RUN ruby-install ruby 2.6.1 && \
    cp /etc/profile.d/chruby.sh /etc/profile.d/chruby-with-ruby-2.6.1.sh && \
    echo "chruby 2.6.1" >> /etc/profile.d/chruby-with-ruby-2.6.1.sh

# Install Bundler
RUN /bin/bash -l -c "gem install bundler --no-document"

# Install Bosh cli
RUN wget https://s3.amazonaws.com/bosh-cli-artifacts/bosh-cli-6.3.1-linux-amd64
RUN chmod +x bosh-cli-*
RUN mv bosh-cli-* /usr/bin/bosh

# Install AWS CLI (used to upload stemcell)
RUN pip install awscli

# receipt generator
RUN cd /tmp && \
    curl -o certify-artifacts -L https://s3.amazonaws.com/bosh-certification-generator-releases/certify-artifacts-linux-amd64 && \
    chmod +x certify-artifacts && \
    mv certify-artifacts /bin/certify-artifacts

RUN useradd non-root-user --shell /bin/bash --create-home
RUN echo "non-root-user ALL=(ALL) NOPASSWD: ALL" >> /etc/sudoers
