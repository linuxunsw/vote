import { defineConfig } from "@hey-api/openapi-ts";

export default defineConfig({
  input: "http://localhost:8888/openapi.yaml",
  output: {
    path: "src/lib/api/internal",
    format: "prettier",
    lint: "eslint",
  },
  plugins: [
    {
      name: "@hey-api/client-fetch",
      runtimeConfigPath: "../config",
    },
    "@hey-api/schemas",
    {
      dates: true,
      name: "@hey-api/transformers",
    },
    {
      enums: "javascript",
      name: "@hey-api/typescript",
    },
    {
      name: "zod",
      compatibilityVersion: 3,
    },
    {
      name: "@hey-api/sdk",
      validator: {
        request: true,
      },
      transformer: true,
    },
  ],
});
