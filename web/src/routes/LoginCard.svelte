<script lang="ts">
  import * as Card from "$lib/components/ui/card";
  import LoginCardOtpStage from "./LoginCardOtpStage.svelte";
  import LoginCardZidStage from "./LoginCardZidStage.svelte";

  let zid = $state<string | null>(null);
</script>

<Card.Root>
  {#if !zid}
    <LoginCardZidStage onsuccess={(providedZid) => (zid = providedZid)} />
  {:else}
    <LoginCardOtpStage {zid} onsuccess={() => {}} oncancel={() => (zid = null)} />
  {/if}
</Card.Root>
