.PHONY: all clean build backend frontend

all: backend frontend admin minio userService
backend: build/backend
frontend: build/frontend
admin: build/admin
minio: build/minio
userService: build/userService


build:
	mkdir -p build

build/backend: build backend/backend
	cp backend/build/backend "$@"

build/minio: build minio/minio
	cp minio/build/minio "$@"

build/userService: build userService/userService
	cp userService/build/userService "$@"

build/frontend: build frontend/frontend
	cp -r frontend/build "$@"

build/admin: build admin/admin
	cp admin/build/admin "$@"

admin/admin:
	make -C "admin" -f "Makefile" all

backend/backend:
	make -C "backend" -f "Makefile" all 

frontend/frontend:
	make -C "frontend" -f "Makefile" all

minio/minio:
	make -C "minio" -f "Makefile" all

userService/userService:
	make -C "userService" -f "Makefile" all