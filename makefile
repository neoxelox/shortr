up:
	docker-compose up -d postgres pgadmin nginx-proxy locust-master locust-worker
	# docker-compose up -d postgres nginx-proxy letsencrypt # Change in production environment
.SILENT: up

stop:
	docker-compose stop
	docker stop $(shell docker ps -a -q)
.SILENT: stop

remove: stop
	docker rm $(shell docker ps -a -q)
.SILENT: remove

prune: stop remove
	docker system prune --force -a
	docker volume prune --force
.SILENT: prune

go/echo: up
	docker-compose up -d --build shortr_go-echo
	docker logs -f shortr_go-echo
.SILENT: go/echo
.PHONY: go/echo
