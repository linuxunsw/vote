import { describe } from "vitest";

describe("Email service", () => {
	// TODO: How to test email service? The below (commented) code is for regular cf workers (i.e. web apis)
	//
	// it("responds with Hello World! (unit style)", async () => {
	// 	const request = new IncomingRequest("http://example.com");
	// 	// Create an empty context to pass to `worker.fetch()`.
	// 	const ctx = createExecutionContext();
	// 	const response = await worker.fetch(request, env, ctx);
	// 	// Wait for all `Promise`s passed to `ctx.waitUntil()` to settle before running test assertions
	// 	await waitOnExecutionContext(ctx);
	// 	expect(await response.text()).toMatchInlineSnapshot(`"Hello World!"`);
	// });
	//
	// it("responds with Hello World! (integration style)", async () => {
	// 	const response = await SELF.fetch("https://example.com");
	// 	expect(await response.text()).toMatchInlineSnapshot(`"Hello World!"`);
	// });
});
