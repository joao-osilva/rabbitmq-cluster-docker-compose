# RabbitMQ cluster with Docker Compose

Creates two clusters with a HAProxy acting as a load balancer.

First cluster has only one RMQ instance. The second has two instances. The haproxy config use SSL termination
and forward to each cluster based on SNI.

This is test case for haproxy SSL termination and SNI routing.

You need to add this line into your /etc/hosts before run the test.

```
127.0.0.1 localhost1.localdomain localhost2.localdomain localhost3.localdomain
```


You need to build the HAProxy image first, just run:
```sh
$ docker build -t haproxy-rabbitmq-cluster:1.7 .
```

Generate the self signed certificates
```sh
$ generate-ssl-key.sh
```

Now run the docker compose file:
```sh
$ docker-compose up -d
```

check if the containers are running:
```sh
$ docker ps
```

Create the cluster by running:
```sh
$ docker exec -ti rabbitmq-node-3 bash -c "rabbitmqctl stop_app"
$ docker exec -ti rabbitmq-node-3 bash -c "rabbitmqctl join_cluster rabbit@rabbitmq-node-2"
$ docker exec -ti rabbitmq-node-3 bash -c "rabbitmqctl start_app"
```

You can run the below script instead
```sh
$ up-cluster.sh
```

Check the cluster status:
```sh
$ docker exec -ti rabbitmq-node-1 bash -c "rabbitmqctl cluster_status"
```

Access HAProxy statistics report at `http://localhost:1936/haproxy?stats` with the credential `haproxy:haproxy`, and the RabbitMQ console at `http://localhost:15672/` with the credential `admin:Admin123`.

To verify that rmq client can send message run the golang program `test-go.go`.
```sh
go build test-go.go
./test-go
```

You can connect to each rmq cluster and verify these messages are received. 
- cluster 0 http://localhost:15672
- cluster 1 http://localhost:15673 (or 15674)
