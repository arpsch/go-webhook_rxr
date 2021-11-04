# go-webhook_rxr
A simple golang based webhook receiver project.

Environmental Variables:
BATACH_INTERVAL in seconds, defaults to 30s if not supplied.
BATCH_SIZE in number, defaults to 3 if not supplied.
ENDPOINT in string targeting the upstream server, defaults to internal test endpoint if not supplied.

How to run:

 docker run -e BATCH_SIZE='5' -e BATCH_INTERVAL='40' ENDPOINT='http://127.0.0.1/logs' -p 9999:9999  <docker_image>:<tag>

