NAME=iot-gps
VERSION=0.0.1
URIDB=postgres://postgres:gtsa4934@34.87.81.32:5432/iot?sslmode=disable
#URIDB=postgres://postgres:gtsa4934@34.87.81.32:5432/iot?sslmode=disable
pg-create:
	## make pg-created NAME="some_awesome_name"
	goose -dir migrations create $(NAME) sql 
pg-status: 
	goose postgres "$(URIDB)" status
pg-up:
	goose -dir migrations postgres "$(URIDB)" up
pg-down:
	goose -dir migrations postgres "$(URIDB)" down
pg-redo:
	goose -dir migrations postgres "$(URIDB)" redo
pg-reset:
	goose -dir migrations postgres "$(URIDB)" reset	
pg-up-by-one: 
	goose -dir migrations postgres "$(URIDB)" up-by-one	
pg-fix: 
	goose -dir migrations postgres "$(URIDB)" fix