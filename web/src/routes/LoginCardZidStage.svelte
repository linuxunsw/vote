<script lang="ts">
  import { generateOtp } from "$lib/api";
  import Button from "$lib/components/ui/button/button.svelte";
  import * as Card from "$lib/components/ui/card";
  import Input from "$lib/components/ui/input/input.svelte";
  import Label from "$lib/components/ui/label/label.svelte";
  import { LoaderCircle } from "@lucide/svelte";
  import * as z from "zod";

  const ZID_VALIDATOR = z.string().regex(/^z\d{7}$/, "Please enter a valid zID (e.g., z1234567).");

  type Props = {
    onsuccess: (zid: string) => void;
  };

  let { onsuccess }: Props = $props();

  let zid = $state("");
  let zidError = $state<string | null>(null);
  let isLoading = $state(false);

  async function handleSubmitZID() {
    zidError = null;
    isLoading = true;

    try {
      ZID_VALIDATOR.parse(zid);
      const { error } = await generateOtp({ body: { zid } });
      if (error) {
        zidError = error.detail ?? "An unknown error occurred. Please try again.";
        return;
      }
      onsuccess(zid);
    } catch (e) {
      if (e instanceof z.ZodError) {
        zidError = e.issues[0].message;
      }
    } finally {
      isLoading = false;
    }
  }
</script>

<Card.Header class="text-center">
  <Card.Title class="text-2xl">Log In</Card.Title>
  <Card.Description>Enter your zID below to receive a one-time password.</Card.Description>
</Card.Header>
<Card.Content>
  <form onsubmit={handleSubmitZID} class="grid gap-4">
    <div class="grid gap-2">
      <Label for="zid">Student ID (zID)</Label>
      <Input
        id="zid"
        placeholder="z1234567"
        bind:value={zid}
        aria-invalid={zidError !== null}
        disabled={isLoading}
      />
      {#if zidError}
        <p class="text-sm text-destructive">{zidError}</p>
      {/if}
    </div>
    <Button type="submit" class="w-full" disabled={isLoading}>
      {#if isLoading}
        <LoaderCircle class="mr-2 h-4 w-4 animate-spin" />
        <span>Sending...</span>
      {:else}
        <span>Continue</span>
      {/if}
    </Button>
  </form>
</Card.Content>
