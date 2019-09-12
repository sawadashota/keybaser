# Chat app needs Keybase app
# Typically, building order is 1) app 2) keybase.
# But app will be changed each time and keybase image can be cache.
# Finally building order is 1) keybase 2) app 3) merge
FROM buildpack-deps:buster-curl

SHELL ["/bin/bash", "-o", "pipefail", "-c"]
RUN \
    apt-get update && \
    apt-get install -y --no-install-recommends \
        fuse \
        libappindicator1 \
        libgconf-2-4 \
        psmisc && \
    # Get and verify Keybase.io's code signing key
    curl -sS https://keybase.io/docs/server_security/code_signing_key.asc | gpg --import &&  \
    gpg --fingerprint 222B85B0F90BE2D24CFEB93F47484E50656D16C7 && \
    # Get, verify and install client package
    curl -sSO https://prerelease.keybase.io/keybase_amd64.deb.sig && \
    curl -sSO https://prerelease.keybase.io/keybase_amd64.deb && \
    gpg --verify keybase_amd64.deb.sig keybase_amd64.deb && \
    apt-get install -y --no-install-recommends ./keybase_amd64.deb && \
    # Create group, user
    groupadd -g 1000 keybase && \
    useradd --create-home -g keybase -u 1000 keybase && \
    # Cleanup
    rm -rf /var/cache/apt/archives/* /var/lib/apt/lists/* && \
    rm keybase_amd64.deb*

USER keybase
WORKDIR /home/keybase
CMD ["bash"]

RUN run_keybase
