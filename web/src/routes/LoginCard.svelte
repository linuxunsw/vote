<script lang="ts">
  import { goto } from "$app/navigation";
  import { resolve } from "$app/paths";
  import type { SubmitOtpResponseBody } from "$lib/api";
  import * as Card from "$lib/components/ui/card";
  import { authState } from "$lib/stores.svelte";
  import LoginCardOtpStage from "./LoginCardOtpStage.svelte";
  import LoginCardZidStage from "./LoginCardZidStage.svelte";

  let zid = $state<string | null>(null);

  function handleLoginSuccess(data: SubmitOtpResponseBody) {
    authState.value = data;
    goto(resolve("/app"));
  }
</script>

<Card.Root>
  {#if !zid}
    <LoginCardZidStage onsuccess={(providedZid) => (zid = providedZid)} />
  {:else}
    <LoginCardOtpStage {zid} onsuccess={handleLoginSuccess} oncancel={() => (zid = null)} />
  {/if}
</Card.Root>
