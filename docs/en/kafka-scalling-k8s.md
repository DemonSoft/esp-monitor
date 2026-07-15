# Kafka Consumer Scaling Architecture in Kubernetes

## 1. Topic Scalability (Topology)

**Per-topic configuration:**  
The number of partitions is configured individually for each topic based on the expected throughput. High-throughput topics (for example, `raw-telemetry`) are divided into multiple partitions (N), while low-throughput topics (such as `alerts-config`) may consist of a single partition.

**Parallelism limit:**  
The number of partitions (N) in a topic defines the maximum number of worker pods that can process messages concurrently within a single Consumer Group.

---

## 2. Worker Lifecycle Management (Who Starts the Workers?)

**Orchestrator (Kubernetes):**  
Kafka itself never creates or manages application pods. The deployment, restart, and maintenance of the desired number of worker replicas are handled by a Kubernetes `Deployment`.

**Autoscaler (KEDA / HPA):**  
A dedicated autoscaling controller (KEDA is recommended) is responsible for dynamically adjusting the number of worker pods.

KEDA periodically queries the Kafka Broker for the `consumer_lag` metric—the number of messages accumulated in topic partitions that have not yet been consumed by the Consumer Group.

Based on predefined scaling rules (for example, *one worker pod for every 500 pending messages*), KEDA dynamically updates the `replicas` field of the Kubernetes `Deployment` within the range of **1...N**, where **N** is the number of partitions in the target topic.

---

## 3. Load Distribution (Who Assigns Partitions?)

**Kafka Group Coordinator:**  
The assignment of partitions to individual worker pods is handled entirely by Kafka's built-in consumer group coordination mechanism (a designated broker acting as the Group Coordinator together with the Consumer Group leader).

**Dynamic Rebalancing:**  

- When KEDA starts a new worker pod, it joins the Consumer Group.
- Kafka automatically triggers a rebalance, redistributing some partitions from existing workers to the newly joined pod.
- If a worker pod fails, Kafka detects the missing heartbeats, initiates another rebalance, and reassigns the released partitions among the remaining active workers.

> 💡 **Running more worker pods than there are partitions in a topic (`Replicas > Partitions`) does not cause system failures. However, the excess pods remain idle and consume cluster resources (CPU and memory) without performing useful work. Therefore, the KEDA autoscaling configuration should set `maxReplicaCount` to a value that is less than or equal to the number of partitions in the target topic.**