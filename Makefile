build:
	docker build -t forum:1.0 .

run:
	docker run -d --name forum-app -p4000:4000 forum:1.0 && echo "\nServer started at http://localhost:4000/"

stop:
	docker stop forum-app

remove:
	docker stop forum-app && docker rm forum-app && docker rmi forum:1.0

.PHONY: build run stop remove
