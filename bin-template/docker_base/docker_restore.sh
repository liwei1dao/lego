# 还原docker 数据卷
tar -zxvf ~/docker-backup/volume.tgz -C ~/docker-backup
docker volume rm mongodata
docker volume rm redisdata
docker run --rm -it -v /var/lib/docker:/docker -v ~/docker-backup/docker/volumes:/volume-backup busybox cp -r /volume-backup/mongodata /docker/volumes
docker run --rm -it -v /var/lib/docker:/docker -v ~/docker-backup/docker/volumes:/volume-backup busybox cp -r /volume-backup/redisdata /docker/volumes