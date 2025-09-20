import { defineConfig } from "@hey-api/openapi-ts";

export default defineConfig({
  input: "http://localhost:8888/openapi.yaml",
  output: {
    path: "src/lib/api",
    format: "prettier",
    lint: "eslint",
  },
  plugins: [
    "@hey-api/schemas",
    {
      dates: true,
      name: "@hey-api/transformers",
    },
    {
      enums: "javascript",
      name: "@hey-api/typescript",
    },
    "zod",
    {
      name: "@hey-api/sdk",
      validator: {
        request: true,
      },
      transformer: true,
    },
  ],
});
