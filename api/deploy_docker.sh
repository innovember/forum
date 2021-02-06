docker image build -f Dockerfile -t "forum" .

echo "----------------------------------------------------------------"

docker images

echo "----------------------------------------------------------------"

docker container run -p 9090:8081 --detach --name forumService forum

echo "----------------------------------------------------------------"

docker ps -a

echo "----------------------------------------------------------------"

docker exec -it forumService /bin/bash
