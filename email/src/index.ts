import { EmailMessage } from "cloudflare:email";
import * as PostalMime from 'postal-mime';
import { createMimeMessage } from 'mimetext';

/*
  “name” (string) - The nominee’s full name
  “contactEmail” (string) (default should be zID email -> up to the client) - The nominee’s contact email
  “statement” (string) - The nominee’s candidate statement
  “roles” (Role array) - An array of the roles (role keys) the nominee is running for
  “discord” (string) - The nominee’s Discord username
  “url” (string) (optional)
*/
export default {
  async email(message, env, ctx) {
    const parser = new PostalMime.default();
    const rawEmail = new Response(message.raw);
    const email = await parser.parse(await rawEmail.arrayBuffer());
    const emailContent = email.html
  },

  async fetch(request, env, ctx) {
    const msg = createMimeMessage();
    msg.setSender({ name: 'Sending email test', addr: 'sender@example.com' });
    msg.setRecipient('recipient@example.com');
    msg.setSubject('An email generated in a worker');
    msg.addMessage({
      contentType: 'text/plain',
      data: `Congratulations, you just sent an email from a worker.`,
    });

    var message = new EmailMessage('sender@example.com', 'recipient@example.com', msg.asRaw());
    await env.EMAIL.send(message);
    return Response.json({ ok: true });
  }
};
