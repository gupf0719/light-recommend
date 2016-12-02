#################################################################
# This docker image build file creates an image that contains
# sqs  project.
#
#                    ##        .
#              ## ## ##       ==
#           ## ## ## ##      ===
#       /""""""""""""""""\___/ ===
#  ~~~ {~~ ~~~~ ~~~ ~~~~ ~~ ~ /  ===- ~~~
#       \______ o          __/
#         \    \        __/
#          \____\______/
#
# Component:    RNTD
# Author:       Gupf <gupf0719@gmail.com>
# Copyright:    (c) 2015-2016 RNTD Ltd. All rights reserved.
#################################################################
#Version 0.0.1
FROM humbleadmin/docker_golang

#go packge
RUN go get github.com/astaxie/beego
RUN go get github.com/beego/bee
RUN go get github.com/astaxie/beego/orm
RUN go get github.com/astaxie/beego/logs
RUN go get github.com/astaxie/beego/swagger
RUN go get github.com/go-sql-driver/mysql

#项目 light-server 配置
ENV Recommend_PATH $GOPATH/src/light-recommend 

RUN mkdir $Recommend_PATH
ADD scripts/docker-recommend/  $Recommend_PATH/
ADD scripts/start.sh /usr/local/

#配置 bee
WORKDIR $GOPATH/src/github.com/beego/bee
RUN go install

#配置 recommend server
WORKDIR $Recommend_PATH/
RUN chmod 777 /usr/local/start.sh

EXPOSE 8080

CMD /usr/local/start.sh
