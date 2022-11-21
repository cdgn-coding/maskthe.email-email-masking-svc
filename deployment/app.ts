import * as kx from "@pulumi/kubernetesx";
import { interpolate } from "@pulumi/pulumi";
import { postgresDatabaseName, postgresSvc } from "./database";
import { rabbitmqSvc } from "./topics";
import { rabbitmqPassword, postgresPassword } from "./config";

const env = {
    "GO_ENVIRONMENT": "production",
    "POSTGRES_DSN": interpolate`postgres://postgres:${postgresPassword}@${postgresSvc.spec.clusterIP}:5432/${postgresDatabaseName}`,
    "RABBITMQ_URL": interpolate`amqp://user:${rabbitmqPassword}@${rabbitmqSvc.spec.clusterIP}:5672/`,
}

const componentName = "email-masking-svc";
const imageName = "email-masking-svc";

const pb = new kx.PodBuilder({
    containers: [
        {
            env,
            name: componentName,
            image: imageName,
            imagePullPolicy: "Never",
            resources: { requests: { cpu: "128m", memory: "256Mi" } },
            ports: { http: 8081 },
            livenessProbe: {
                httpGet: {
                    path: "/health",
                    port: 8081,
                },
            },
        },
    ],
});

const deployment = new kx.Deployment(componentName, {
    spec: pb.asDeploymentSpec({ replicas: 1 }),
});

export const appService = deployment.createService({
    type: kx.types.ServiceType.ClusterIP,
});
