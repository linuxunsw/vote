import * as PostalMime from 'postal-mime';
import { EmailMessage } from "cloudflare:email";
import { createMimeMessage } from 'mimetext';

export default {
  async email(message, env, ctx) {
    const msg = createMimeMessage();
    msg.setSender({ name: 'Thank you for your contact', addr: 'sender@example.com' });
    msg.setRecipient(message.from);
    msg.setHeader('In-Reply-To', message.headers.get('Message-ID'));
    msg.setSubject('An email generated in a worker');
    msg.addMessage({
      contentType: 'text/plain',
      data: `This is an automated reply. We received your email.`,
    });

    const replyMessage = new EmailMessage('sender@example.com', message.from, msg.asRaw());

    await message.reply(replyMessage);
  },
};
