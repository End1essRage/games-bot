#defaults
port ?= 8080
env = ENV_POD

#для докера
d_connect:
	docker exec -it ${container} /bin/bash

d_run:
	docker run --rm -p ${port}:8080 -e ENV=$(env) -e TOKEN=$(t) end1essrage/games-bot

d_build: 
	docker build -t end1essrage/games-bot .

#для подмена
p_connect:
	podman exec -it ${container} /bin/bash
	
p_run:
	podman run -p ${port}:8080 -e ENV=$(env) -e TOKEN=$(t) end1essrage/games-bot

p_build: 
	podman build -t end1essrage/games-bot .