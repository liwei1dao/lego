# 备份docker 数据卷
docker run --rm -it -v ~/docker-backup:/backup -v /var/lib/docker:/docker busybox tar cfz /backup/volume.tgz /docker/volumes/