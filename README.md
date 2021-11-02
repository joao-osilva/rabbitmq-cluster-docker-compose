# HAProxy and RabbitMQ cluster with Docker Compose

Creates a 3 node RabbitMQ cluster with a HAProxy acting as a load balancer.

You need to build the HAProxy image first, just run:
```sh
docker build -t roofimon/haproxy-rabbitmq-cluster:2.3 .
```

Now run the docker compose file:
```sh
docker-compose up -d
```

check if the containers are running:
```sh
docker ps
```

At this stage haproxy can connect to all 3 nodes of rabbitmq but those 3 nodes still working seperately. As a result, we need to do one more step to form a cluster.
We can see status of this stage by run the cluster status command:
```sh
docker exec -ti rabbitmq-node-1 bash -c "rabbitmqctl cluster_status"
```
Result will be like the detail below. There is only one node, rabbitmq-node-1, in cluster
```sh
Cluster status of node rabbit@rabbitmq-node-1 ...
Basics

Cluster name: rabbit@rabbitmq-node-1

Disk Nodes

rabbit@rabbitmq-node-1

Running Nodes

rabbit@rabbitmq-node-1
```

Create the cluster by running all command below. Basically it will tell other 2 nodes to join rabbit-node-1 and form cluster:
```sh
docker exec -ti rabbitmq-node-2 bash -c "rabbitmqctl stop_app"
docker exec -ti rabbitmq-node-2 bash -c "rabbitmqctl join_cluster rabbit@rabbitmq-node-1"
docker exec -ti rabbitmq-node-2 bash -c "rabbitmqctl start_app"

docker exec -ti rabbitmq-node-3 bash -c "rabbitmqctl stop_app"
docker exec -ti rabbitmq-node-3 bash -c "rabbitmqctl join_cluster rabbit@rabbitmq-node-1"
docker exec -ti rabbitmq-node-3 bash -c "rabbitmqctl start_app"
```
Or we can just run shell script that pack all commands above into one file:
```sh
chmod 777 create-cluster.sh
./create-cluster.sh
```

Check the cluster status:
```sh
docker exec -ti rabbitmq-node-1 bash -c "rabbitmqctl cluster_status"
```
We now can see that there 3 nodes running under the cluster
```sh
Cluster status of node rabbit@rabbitmq-node-1 ...
Basics

Cluster name: rabbit@rabbitmq-node-1

Disk Nodes

rabbit@rabbitmq-node-1
rabbit@rabbitmq-node-2
rabbit@rabbitmq-node-3

Running Nodes

rabbit@rabbitmq-node-1
rabbit@rabbitmq-node-2
rabbit@rabbitmq-node-3
```

Options for more resilience 

Declares a policy which matches the queues whose names begin with "two." are mirrored to any two nodes in the cluster, with automatic synchronisation:
```	
docker exec -ti rabbitmq-node-3 bash -c 'rabbitmqctl set_policy ha-two "^two\." "{\"ha-mode\":\"exactly\",\"ha-params\":2,\"ha-sync-mode\":\"automatic\"}"'
```

Declares a policy which matches the queues whose names begin with "test_" and configures mirroring to all nodes in the cluster:
```sh
docker exec -ti rabbitmq-node-3 bash -c 'rabbitmqctl set_policy ha-all "^test\_" "{\"ha-mode\":\"all\"}"'
```

Access HAProxy statistics report at 
```
http://localhost:1936/haproxy?stats` with the credential `haproxy:haproxy`
```
RabbitMQ console at 
```
http://localhost:15672/` with the credential `admin:Admin@123`.
```
