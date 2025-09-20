import type { PageLoad } from "./$types";

export const load: PageLoad = () => {
  // TODO: api stuffs

  // dummy data
  return {
    state: {
      election_id: "asdf",
      state: "NO_ELECTION",
    },
  };
};
