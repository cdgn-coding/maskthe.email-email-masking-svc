import * as k8s from "@pulumi/kubernetes";
import * as kx from "@pulumi/kubernetesx";
import { rabbitmqPassword } from "./config";

export const rabbitmqSecret = new k8s.core.v1.Secret("rabbitmq-secret", {
    stringData: {
        "user": "user",
        "password": rabbitmqPassword
    }
});

const volume = new k8s.core.v1.PersistentVolume("rabbitmq-pv", {
    metadata: {
        name: "rabbitmq-pv",
    },
    spec: {
        accessModes: ["ReadWriteOnce"],
        capacity: {
            storage: "5Gi",
        },
        hostPath: {
            path: "/data/pv-rabbitmq"
        }
    },
});

const volumeClaim = new k8s.core.v1.PersistentVolumeClaim("rabbitmq-pvc", {
    metadata: {
        name: "rabbitmq-pvc",
    },
    spec: {
        accessModes: ["ReadWriteOnce"],
        resources: {
            requests: {
                storage: "2Gi",
            }
        },
    },
});

const pb = new kx.PodBuilder({
    volumes: [
        {
            name: "rabbitmq-volume-storage",
            persistentVolumeClaim: {
                claimName: volumeClaim.metadata.name,
            }
        }
    ],
    containers: [
        {
            name: "rabbitmq",
            image: "docker.io/bitnami/rabbitmq:latest",
            imagePullPolicy: "IfNotPresent",
            resources: { requests: { cpu: "128m", memory: "256Mi" } },
            env: [
                {
                    name: "RABBITMQ_USERNAME",
                    valueFrom: {
                        secretKeyRef: {
                            name: rabbitmqSecret.metadata.name,
                            key: "user"
                        }
                    }
                },
                {
                    name: "RABBITMQ_PASSWORD",
                    valueFrom: {
                        secretKeyRef: {
                            name: rabbitmqSecret.metadata.name,
                            key: "password"
                        }
                    }
                },
            ],
            ports: [
                {
                    name: "amqp",
                    protocol: "TCP",
                    containerPort: 5672,
                },
                {
                    name: "http",
                    protocol: "TCP",
                    containerPort: 15672,
                }
            ],
            volumeMounts: [
                {
                    name: "rabbitmq-volume-storage",
                    mountPath: "/opt/bitnami/rabbitmq/etc/rabbitmq/"
                },
            ]
        }
    ],
});

// Deploy postgres as a StatefulSet.
export const rabbitmqDeployment = new kx.Deployment("rabbitmq-deployment", {
    spec: pb.asDeploymentSpec({ replicas: 1 }),
});

export const rabbitmqSvc = rabbitmqDeployment.createService({
    type: kx.types.ServiceType.ClusterIP,
});

export const rabbitmqClusterIP = rabbitmqSvc.spec.clusterIP;