up:
	docker-compose up -d postgres pgadmin
.SILENT: up

go/echo: up
	docker-compose up -d --build shortr_go-echo
	docker logs -f shortr_go-echo
.SILENT: go/echo
.PHONY: go/echo

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