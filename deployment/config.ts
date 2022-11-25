import * as pulumi from "@pulumi/pulumi";

const config = new pulumi.Config();

const clusterSetup = new pulumi.StackReference(
    config.require("cluster-setup")
);

export const postgresDsn = clusterSetup.requireOutput("postgresDsn");
export const rabbitmqEndpoint = clusterSetup.requireOutput("rabbitmqEndpoint");