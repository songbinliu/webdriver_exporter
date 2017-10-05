FROM ubuntu:16.04

#1. install google chrome
RUN sh -c 'echo "deb [arch=amd64] http://dl.google.com/linux/chrome/deb/ stable main" >> /etc/apt/sources.list.d/google.list'
RUN apt-get update && apt-get install -y --no-install-recommends apt-utils
RUN apt-get install  xvfb -y
RUN apt-get install google-chrome-stable  -y --allow-unauthenticated

# RUN apt-get install chromium-chromedriver -y
# RUN ln -s /usr/lib/chromium-browser/chromedriver /usr/bin/chromedriver 
COPY ./chromedriver /usr/local/bin/chromedriver
RUN chmod 0755 /usr/local/bin/chromedriver

#2. install prometheus
RUN mkdir /etc/prometheus
COPY ./prometheus/prometheus                             /bin/prometheus
COPY ./prometheus/promtool                               /bin/promtool
COPY ./prometheus/prometheus.yml                         /etc/prometheus/prometheus.yml
COPY ./prometheus/console_libraries/                     /etc/prometheus/
COPY ./prometheus/consoles/                              /etc/prometheus/

RUN chmod 0755 /bin/prometheus
EXPOSE 9090
VOLUME     [ "/prometheus" ]

#3. install wedriver_exporter
COPY ./webdriver.linux /usr/bin/webdriver
RUN chmod 0755 /usr/bin/webdriver

#4. the script to provide parameters
COPY ./run.sh /usr/bin/run_webdriver.sh
RUN chmod 0755 /usr/bin/run_webdriver.sh

#4 Port
EXPOSE 9156

ENTRYPOINT ["/usr/bin/run_webdriver.sh"]
