import * as kx from "@pulumi/kubernetesx";
import * as k8s from "@pulumi/kubernetes";
import { createService } from "./utils";
import { postgresDsn, rabbitmqEndpoint } from "./config";


const env = {
    "GO_ENVIRONMENT": "production",
    "POSTGRES_DSN": postgresDsn,
    "RABBITMQ_URL": rabbitmqEndpoint,
}

const componentName = "email-masking-svc";
const imageName = "email-masking-svc:v1";

const pb = new kx.PodBuilder({
    containers: [
        {
            env,
            name: componentName,
            image: imageName,
            imagePullPolicy: "IfNotPresent",
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

export const appService = createService({
    name: componentName,
    serviceSpecs: {
        type: kx.types.ServiceType.ClusterIP,
        ports: [{
            protocol: "TCP",
            name: "web",
            port: 8080,
            targetPort: 8081
        }],
        selector: {
            app: componentName,
        },
    },
    metadata: {
        name: componentName,
    }
}, deployment);

export const appEndpoint = `${componentName}.default.svc.cluster.local`;
