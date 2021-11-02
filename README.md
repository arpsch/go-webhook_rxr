# go-webhook_rxr
A simple golang based webhook receiver project.

How to run:

 docker run -e BATCH_SIZE='5' -e BATCH_INTERVAL='40' ENDPOINT='http://127.0.0.1/logs' -p 9999:9999  <docker_image>:<tag>
