import { goto } from "$app/navigation";
import { getElectionState } from "$lib/api";
import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";

export const prerender = false;
export const ssr = false;

export const load: PageLoad = async ({ fetch }) => {
  const { data, error: errorData } = await getElectionState({ fetch });

  if (errorData) {
    const errorCode = Number(errorData.status);
    if (errorCode === 401) {
      goto("/?error=session_expired");
    } else if (errorCode === 403) {
      goto("/?error=unauthorized");
    } else {
      error(errorCode, errorData.detail);
    }
  }

  if (!data) {
    // This should be unreachable - it's mainly to appease typescript
    error(500, "An error occurred.");
  }

  // dummy data
  return {
    state: data,
  };
};
