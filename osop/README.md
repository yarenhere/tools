
```azure
docker run --rm -it  -p 9000:9000 -p 9001:9001 --name minio -v ~/minio-server/:/data minio/minio server --address ":9000" --console-address ":9001" /data
```
启动minio作为对象存储后端进行测试



