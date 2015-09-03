FROM ruby:2.2

ENV http_proxy=http://proxy.***REMOVED***:8080 \
    https_proxy=http://proxy.***REMOVED***:8080 \
    no_proxy=***REMOVED***,localhost,127.0.0.1

# throw errors if Gemfile has been modified since Gemfile.lock
RUN bundle config --global frozen 1

RUN mkdir -p /app
WORKDIR /app

COPY Gemfile /app/
COPY Gemfile.lock /app/
RUN bundle install

RUN git config --global user.email middleman@middleman.local && \
    git config --global user.name middleman 

RUN echo "StrictHostKeyChecking no" > /etc/ssh/ssh_config

COPY deploy.sh /deploy.sh



