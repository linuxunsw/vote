import { defineConfig } from "@hey-api/openapi-ts";

export default defineConfig({
  input: "http://localhost:8888/openapi.yaml",
  output: "./src/lib/api",
});
