import * as pulumi from "@pulumi/pulumi";
import * as k8s from "@pulumi/kubernetes";

const config = new pulumi.Config();

const clusterSetup = new pulumi.StackReference(
    config.require("cluster-setup")
);

export const postgresDsn = clusterSetup.requireOutput("postgresDsn");
export const rabbitmqEndpoint = clusterSetup.requireOutput("rabbitmqEndpoint");

const kubernetesStack = new pulumi.StackReference(config.require("kubernetesStack"));

const kubeconfig = kubernetesStack.getOutput("kubeconfig");
export const k8sProvider = new k8s.Provider("k8s", { kubeconfig: kubeconfig });
export const baseOptions = { provider: k8sProvider };

export const registryUrl = kubernetesStack.getOutput("registryUrl");
export const registryEndpoint = kubernetesStack.getOutput("registryEndpoint");
export const registryUser = kubernetesStack.getOutput("registryUser");
export const registryPassword = kubernetesStack.getOutput("registryPassword");
export const registryName = kubernetesStack.getOutput("registryName");

export const currentStack = pulumi.getStack();
export const imageName = config.require("imageName")