## INSTRUCTION

## DATABASE INSTRUCTION:

Database docker instruction:

docker pull postgres:15.1-alpine

docker run --name postgres15.1 -e POSTGRES_DB=tmdt-db -e POSTGRES_USER=root -e POSTGRES_PASSWORD=1 -p 8001:5432 -d postgres:15.1-alpine


docker run -d --name elasticsearch -p 9200:9200 -e discovery.type=single-node -v elasticsearch:/usr/share/elasticsearch/data docker.elastic.co/elasticsearch/elasticsearch:8.6.1

docker run --name elasticsearch --net elastic -p 9200:9200 -e discovery.type=single-node -e ES_JAVA_OPTS="-Xms1g -Xmx1g" -e xpack.security.enabled=false -it docker.elastic.co/elasticsearch/elasticsearch:8.6.1