FROM registry.access.redhat.com/redhat-openjdk-18/openjdk18-openshift

MAINTAINER milkliver
#ARG uid=0
#ARG gid=0
USER 0

#==========================Add resources==========================

RUN mkdir -p /etc/opt/scdf
ADD ./resources/* /etc/opt/scdf/

#===========================install rpms==========================
RUN mkdir /rpms
WORKDIR /rpms
ADD ./rpms /rpms
RUN rpm -ivh --nodigest --nofiledigest /rpms/*

RUN java -version



#===================add scdf executor and jobs===================
RUN mkdir /testfiles
WORKDIR /testfiles

ADD ./*.jar /testfiles/
RUN chmod 777 -Rf /testfiles

RUN mkdir /configs
ADD ./execution.properties /configs/execution.properties
RUN chmod 777 -Rf /configs/execution.properties



#============================run scdf===========================
USER 1001

#ENTRYPOINT ["/bin/sh", "-c" , "echo \"192.168.50.135 broker\" >> /etc/hosts&& /bin/java -jar -Dspring.config.location=/configs/execution.properties /testfiles/executor.jar"]
ENTRYPOINT ["/bin/java","-jar","-Dspring.config.location=/configs/execution.properties","/testfiles/executor.jar"]
#CMD ["/bin/java","-jar","-Dspring.config.location=/configs/execution.properties","/testfiles/scdf-task01.jar"]


# For Test
#CMD ["tail","-f","/dev/null"]
