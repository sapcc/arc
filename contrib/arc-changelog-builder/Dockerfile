FROM alpine:3.6

ENV GITHUB_CHANGELOG_GENERATOR_VERSION "1.14.3"

RUN apk --no-cache add ruby ruby-json libstdc++ tzdata bash ca-certificates git curl
RUN curl http://aia.pki.co.sap.com/aia/SAP%20Global%20Root%20CA.crt | \
    tr -d '\r' > /usr/local/share/ca-certificates/SAP_Global_Root_CA.crt && \
    update-ca-certificates

RUN echo 'gem: --no-document' > /etc/gemrc
RUN git clone https://github.com/ArtieReus/github-changelog-generator.git -b bisect2
WORKDIR /github-changelog-generator
RUN gem build github_changelog_generator.gemspec
RUN gem install github_changelog_generator-1.15.0.pre.rc.gem
WORKDIR /

COPY .github_changelog_generator /

CMD github_changelog_generator --token $GITHUB_TOKEN --github-api $GITHUB_API
