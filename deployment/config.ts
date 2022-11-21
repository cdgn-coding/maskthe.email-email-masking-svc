import * as pulumi from "@pulumi/pulumi";

const config = new pulumi.Config();

export const rabbitmqPassword = config.requireSecret("rabbitmqPassword");
export const postgresPassword = config.requireSecret("postgresPassword");
