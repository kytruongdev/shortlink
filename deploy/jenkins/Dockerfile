# Jenkins with the Docker CLI + Compose plugin (to deploy the stack) and a
# preconfigured setup: plugins baked in, setup wizard skipped, job created via
# Configuration-as-Code. No manual clicking to bootstrap.
FROM jenkins/jenkins:lts-jdk17

USER root
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates curl && \
    install -m 0755 -d /etc/apt/keyrings && \
    curl -fsSL https://download.docker.com/linux/debian/gpg -o /etc/apt/keyrings/docker.asc && \
    chmod a+r /etc/apt/keyrings/docker.asc && \
    echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/debian $(. /etc/os-release && echo $VERSION_CODENAME) stable" > /etc/apt/sources.list.d/docker.list && \
    apt-get update && apt-get install -y --no-install-recommends docker-ce-cli docker-compose-plugin && \
    rm -rf /var/lib/apt/lists/*

COPY casc.yaml /var/jenkins_config/casc.yaml
ENV CASC_JENKINS_CONFIG=/var/jenkins_config/casc.yaml
ENV JAVA_OPTS="-Djenkins.install.runSetupWizard=false"

USER jenkins
COPY plugins.txt /usr/share/jenkins/ref/plugins.txt
RUN jenkins-plugin-cli -f /usr/share/jenkins/ref/plugins.txt
