import * as kx from "@pulumi/kubernetesx";
import { postgresConnectionDsn } from "./database";
import { rabbitmqUrl } from "./topics";

const env = {
    "GO_ENVIRONMENT": "production",
    "POSTGRES_DSN": postgresConnectionDsn,
    "RABBITMQ_URL": rabbitmqUrl,
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
    ports: [{ port: 8080, targetPort: 8081 }],
});
