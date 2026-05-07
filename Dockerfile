
FROM  ubuntu:latest
EXPOSE 7890
WORKDIR /service
COPY ./deploy/bin/server_go /service/server_go 
COPY ./deploy/bin/config/ /service/config/ 
COPY ./deploy/bin/Resources /service/Resources
 
RUN apt update
RUN apt -y upgrade
RUN apt-get install -y ffmpeg

RUN apt-get install -y openjdk-8-jdk
RUN apt install -y libreoffice
RUN apt install -y p7zip-full p7zip-rar
CMD ["/service/server_go"]


# FROM ubuntu:latest
# RUN apt-get update && \
#     apt-get install -y ffmpeg
# # 设置容器启动时执行的命令
# CMD [ "ffmpeg", "-version" ]







