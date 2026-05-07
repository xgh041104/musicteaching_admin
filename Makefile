ARCH := $(if $(ARCH),$(ARCH),amd64)

build:
	rm -rf ./deploy/bin || true
	mkdir ./deploy/bin || true
	mkdir ./deploy/bin/Resources || true
	cd  server_go/cmd/server && GOOS=linux  GOARCH=$(ARCH) go build -o ../../../deploy/bin/server_go
	cp  -r ./server_go/config/ ./deploy/bin
	#cp -r ./server_go/Resources/SlideImg ./deploy/bin/Resources/SlideImg

clean:
	docker image rm mysql nginx server  || true
	rm ./mysql.tar || true
	rm ./nginx.tar || true
	rm ./server_go.tar || true
	rm ./deploy/mysql.tar || true
	rm ./deploy/nginx.tar || true
	rm ./deploy/server_go.tar || true
	sudo rm -rf ./deploy/mysql/data/* ||true
	rm ./deploy/mysql/init.d/* || true
	rm -rf ./deploy/website/manage/* || true
	rm -rf ./deploy/website/student/* || true

dep:
	# make clean
	#cp ./db/subjectcourse.sql ./deploy/mysql/init.d/init_01.sql
	cp -r ./web_react/dist/* ./deploy/website/manage/
	#cp -r ./react_web_student/dist/* ./deploy/website/student/


	 #docker pull --platform $(ARCH) mysql:latest
	 #docker pull --platform $(ARCH) nginx:latest
	 #docker save -o mysql.tar mysql:latest
	 #docker save -o nginx.tar nginx:latest

	
	#cp mysql.tar ./deploy/
	#cp nginx.tar ./deploy/

dockerx:
	#docker build -t server_go:latest .
	docker save -o server_go.tar server_go:latest
	cp server_go.tar ./deploy/

docker:
	docker build -t server_go:latest -f Dockerfile .
 


