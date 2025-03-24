#!/bin/sh

# Wait for Elasticsearch to be ready
until curl -s http://localhost:9200/_cluster/health | grep -q 'status.*green\|status.*yellow'; do
    echo 'Waiting for Elasticsearch to be ready...'
    sleep 1
done

# Create the catalog index if it doesn't exist
if curl -s -f "http://localhost:9200/catalog" > /dev/null; then
    echo "Index 'catalog' already exists"
else
    echo "Creating 'catalog' index..."
    curl -X PUT "http://localhost:9200/catalog" \
         -H "Content-Type: application/json" \
         -d @/usr/share/elasticsearch/init/init.json
    echo "\nIndex 'catalog' created successfully"
fi 