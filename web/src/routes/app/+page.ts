import { goto } from "$app/navigation";
import { resolve } from "$app/paths";
import { getNomination } from "$lib/api";
import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";

export const prerender = false;
export const ssr = false;

export const load: PageLoad = async () => {
  const { data, error: errorData } = await getNomination();

  console.log(data, errorData);
  if (errorData) {
    const errorCode = Number(errorData.status);
    if (errorCode === 401) {
      // logged out or session expired, redirect to login
      // TODO: add some user feedback? like message saying "session expired" or smth
      goto(resolve("/"));
      return;
    } else if (errorCode !== 404) {
      // 404 means no nomination, which isn't a fatal error for this client
      error(errorCode, errorData.detail);
    }
  }

  // the idea is that whenever the user updates any of this stuff, the invalidateAll() function will be called, causing an update

  return {
    nomination: data,
    ballot: {
      // TODO
    },
  };
};
