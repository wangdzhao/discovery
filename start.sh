docker build -t discovery:latest .
docker run -d -p 9999:8888 --name discovery_demo discovery:latest
