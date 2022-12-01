import * as pulumi from "@pulumi/pulumi";
import * as docker from "@pulumi/docker";
import {
  registryUrl,
  registryEndpoint,
  registryUser,
  registryPassword,
  currentStack,
  imageName,
} from "./config";

const getImageName = () => {
  if (currentStack === "dev") {
    return imageName;
  }

  const image = new docker.Image("appImage", {
    build: "../",
    localImageName: `${imageName}:v1`,
    imageName: pulumi.interpolate`${registryEndpoint}/${imageName}:v1`,
    registry: {
        server: registryUrl,
        username: registryUser,
        password: registryPassword,
    },
  });

  return image.imageName;
}

export const fullImageName = getImageName();