FROM elasticsearch:7.17.9

# Install curl
USER root
RUN apt-get update && apt-get install -y curl
USER elasticsearch

# Copy initialization files
COPY init.json /usr/share/elasticsearch/init/
COPY init.sh /usr/share/elasticsearch/init/

# Make the script executable
USER root
RUN chmod +x /usr/share/elasticsearch/init/init.sh
USER elasticsearch

# Add an entrypoint script that runs our init script after Elasticsearch starts
COPY --chown=elasticsearch:elasticsearch docker-entrypoint.sh /usr/local/bin/custom-docker-entrypoint.sh
USER root
RUN chmod +x /usr/local/bin/custom-docker-entrypoint.sh
USER elasticsearch

ENTRYPOINT ["/usr/local/bin/custom-docker-entrypoint.sh"] 