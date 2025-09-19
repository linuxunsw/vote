import type { PageLoad } from "./$types";

export const load: PageLoad = () => {
  // TODO: check if user is logged in and if not, navigate to /
  // TODO: get user's nomination and voting status

  // the idea is that whenever the user updates any of this stuff, the invalidateAll() function will be called, causing an update

  return {
    nomination: {
      candidate_name: "John Doe",
      candidate_statement: "I am running for president because...",
      contact_email: "john@example.com",
      discord_username: "johndoe",
      executive_roles: ["president", "secretary"],
      url: "https://johndoe.com",
    },
    ballot: {
      // TODO
    },
  };
};
