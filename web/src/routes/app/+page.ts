import { goto } from "$app/navigation";
import { getBallot, getNomination } from "$lib/api";
import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";

export const prerender = false;
export const ssr = false;

export const load: PageLoad = async ({ fetch }) => {
  const [nominationRes, ballotRes] = await Promise.all([
    getNomination({ fetch }),
    getBallot({ fetch }),
  ]);

  const { data: nomination, error: errorData } = nominationRes;
  const { data: ballot } = ballotRes;

  if (errorData) {
    const errorCode = Number(errorData.status);
    if (errorCode === 401) {
      // logged out or session expired, redirect to login
      goto("/?error=session_expired");
      return;
    } else if (errorCode !== 404 && errorCode !== 400) {
      // 404 means no nomination and 400 means no election, which aren't fatal errors for this client
      error(errorCode, errorData.detail);
    }
  }

  // the idea is that whenever the user updates any of this stuff, the invalidateAll() function will be called, causing an update

  return {
    nomination,
    ballot,
  };
};
