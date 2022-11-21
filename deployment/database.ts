import * as k8s from "@pulumi/kubernetes";
import * as kx from "@pulumi/kubernetesx";
import { postgresPassword} from "./config";
import { createService } from "./utils";
import {interpolate} from "@pulumi/pulumi";

export const postgresDatabaseName = "emails";

export const postgresSecret = new k8s.core.v1.Secret("postgres", {
    stringData: {
        "user": "postgres",
        "password": postgresPassword
    }
});

const volume = new k8s.core.v1.PersistentVolume("postgres-pv", {
    metadata: {
        name: "postgres-pv",
    },
    spec: {
        accessModes: ["ReadWriteOnce"],
        capacity: {
            storage: "5Gi",
        },
        hostPath: {
            path: "/data/pv-postgres"
        }
    },
});

export const volumeClaim = new k8s.core.v1.PersistentVolumeClaim("postgres-pvc", {
    metadata: {
        name: "postgres-pvc",
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
            name: "postgres-volume-storage",
            persistentVolumeClaim: {
                claimName: volumeClaim.metadata.name,
            }
        }
    ],
    containers: [
        {
            name: "postgres",
            image: "docker.io/postgres:13.9",
            imagePullPolicy: "IfNotPresent",
            resources: { requests: { cpu: "128m", memory: "256Mi" } },
            env: [
                {
                    name: "POSTGRES_USER",
                    valueFrom: {
                        secretKeyRef: {
                            name: postgresSecret.metadata.name,
                            key: "user"
                        }
                    }
                },
                { name: "POSTGRES_DB", value: postgresDatabaseName },
                {
                    name: "POSTGRES_PASSWORD",
                    valueFrom: {
                        secretKeyRef: {
                            name: postgresSecret.metadata.name,
                            key: "password"
                        }
                    }
                },
            ],
            ports: [
                {
                    name: "postgres",
                    protocol: "TCP",
                    containerPort: 5432,
                }
            ],
            livenessProbe: {
                exec: {
                    command: ["sh", "-c", `PGPASSWORD=$POSTGRES_PASSWORD psql -w -U "postgres" -d "postgres"  -h 127.0.0.1 -c "SELECT 1"`],
                },
                initialDelaySeconds: 120,
                periodSeconds: 10,
                timeoutSeconds: 1,
                successThreshold: 1,
                failureThreshold: 3
            },
            readinessProbe: {
                exec: {
                    command: ["sh", "-c", `PGPASSWORD=$POSTGRES_PASSWORD psql -w -U "postgres" -d "postgres"  -h 127.0.0.1 -c "SELECT 1"`],
                },
                initialDelaySeconds: 30,
                periodSeconds: 10,
                timeoutSeconds: 1,
                successThreshold: 1,
                failureThreshold: 3
            },
            volumeMounts: [
                {
                    name: "postgres-volume-storage",
                    mountPath: "/var/lib/postgresql"
                },
            ]
        }
    ],
});

// Deploy postgres as a StatefulSet.
export const postgresDeployment = new kx.Deployment("postgres-deployment", {
    spec: pb.asDeploymentSpec({ replicas: 1 }),
});

const serviceName = "email-masking-postgres";

export const postgresSvc = createService(
    {
        name: serviceName,
        serviceSpecs: { type: kx.types.ServiceType.ClusterIP },
        metadata: {
            name: serviceName,
        },
    },
    postgresDeployment
);

export const postgresClusterIP = postgresSvc.spec.clusterIP;

export const postgresConnectionDsn = interpolate`postgres://postgres:${postgresPassword}@${serviceName}.default.svc.cluster.local:5432/${postgresDatabaseName}`;