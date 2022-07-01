FROM docker.io/milkliver/javamvnbuilder
MAINTAINER milkliver
#ARG uid=0
#ARG gid=0

VOLUME /sys/fs/cgroup

RUN mkdir /tmp/src
COPY ./ /tmp/src/
RUN chown -R 1001:1001 /tmp/src

USER 1001

WORKDIR /opt/app-root/run

# Assemble script sourced from builder image based on user input or image metadata.
# If this file does not exist in the image, the build will fail.
RUN ["/opt/app-root/run/assemble"]
# Run script sourced from builder image based on user input or image metadata.
# If this file does not exist in the image, the build will fail.
ENTRYPOINT ["/opt/app-root/run/run"]
#CMD ["/opt/app-root/run/run"]

