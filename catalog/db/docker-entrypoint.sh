#!/bin/bash

# Start Elasticsearch in the background
/usr/local/bin/docker-entrypoint.sh eswrapper &

# Wait for Elasticsearch to be ready
until curl -s http://localhost:9200/_cluster/health | grep -q 'status.*green\|status.*yellow'; do
    echo 'Waiting for Elasticsearch to be ready...'
    sleep 2
done

# Run our initialization script
/usr/share/elasticsearch/init/init.sh

# Keep the container running
wait 