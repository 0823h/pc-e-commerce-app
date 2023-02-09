## INSTRUCTION

## DATABASE INSTRUCTION:

Database docker instruction:

docker pull postgres:15.1-alpine

docker run --name postgres15.1 -e POSTGRES_DB=tmdt-db -e POSTGRES_USER=root -e POSTGRES_PASSWORD=1 -p 8001:5432 -d postgres:15.1-alpine
