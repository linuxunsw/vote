import type { CreateClientConfig } from "./internal/client.gen";

export const createClientConfig: CreateClientConfig = (config) => ({
  ...config,
  baseUrl: "/",
});
