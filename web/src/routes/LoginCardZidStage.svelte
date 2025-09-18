<script lang="ts">
  import Button from "$lib/components/ui/button/button.svelte";
  import * as Card from "$lib/components/ui/card";
  import Input from "$lib/components/ui/input/input.svelte";
  import { ArrowRight } from "@lucide/svelte";
  import * as z from "zod";

  const ZID_VALIDATOR = z.string().regex(/^z\d{7}$/);

  type Props = {
    onsuccess: (zid: string) => void;
  };

  let { onsuccess }: Props = $props();

  let zid = $state("");
  let zidError = $state<string | null>(null);

  async function handleSubmitZID() {
    console.log("hit");
    zidError = null;
    try {
      ZID_VALIDATOR.parse(zid);
      // TODO: Actual API stuff
      onsuccess(zid);
    } catch (e) {
      if (e instanceof z.ZodError) {
        zidError = "Invalid zID entered.";
      }
    }
  }
</script>

<Card.Header>
  <Card.Title>Login with zID</Card.Title>
</Card.Header>
<Card.Content>
  <form onsubmit={handleSubmitZID} class="flex w-full flex-col gap-1.5">
    <Input placeholder="zXXXXXXX" bind:value={zid} aria-invalid={zidError !== null} />
    {#if zidError !== null}
      <p class="text-sm text-muted-foreground">{zidError}</p>
    {/if}
    <Button class="gap-1" type="submit">
      <span>Log In</span>
      <ArrowRight />
    </Button>
  </form>
</Card.Content>
