FROM ubuntu:latest
# Set working directory
WORKDIR /opt/app

# Pull the distribution over
COPY /tmp/gwstlg-1.0.2.tar.gz /opt/app/gwstlg-1.0.2.tar.gz
RUN mkdir -p /opt/app
RUN mkdir -p /opt/app/etc
RUN mkdir -p /opt/app/logs
RUN cd /opt/app && tar -zxf /tmp/gwstlg-1.0.2.tar.gz

# Run the server
CMD /opt/app/gwstlg-1.0.2/bin/grow-with-stl-go --loglevel 6
