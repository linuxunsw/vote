import { EmailMessage } from "cloudflare:email";
import { createMimeMessage } from "mimetext";

export default {
  async email(message, env, ctx) {
    const allowList = ["nsdj.sharma@gmail.com"];
    if (!allowList.includes(message.from)) {
      message.setReject("Address not allowed");
      return;
    }

    const reply = createMimeMessage();
    reply.setHeader("In-Reply-To", message.headers.get("Message-ID")!);
    reply.setSender({ name: "Thank you for your contact", addr: "election@linuxunsw.org" });
    reply.setRecipient(message.from);
    reply.setSubject("Email test successful");
    reply.addMessage({
      contentType: "text/plain",
      data: "Pong",
    });

    const replyMessage = new EmailMessage("<SENDER>@example.com", message.from, reply.asRaw());

    await message.reply(replyMessage);
  },
} satisfies ExportedHandler<Env>;
