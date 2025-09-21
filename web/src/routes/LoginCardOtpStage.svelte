<script lang="ts">
  import { submitOtp, type SubmitOtpResponseBody } from "$lib/api";
  import Button from "$lib/components/ui/button/button.svelte";
  import * as Card from "$lib/components/ui/card";
  import * as InputOTP from "$lib/components/ui/input-otp";
  import { ArrowRight } from "@lucide/svelte";
  import * as z from "zod";

  const OTP_VALIDATOR = z.string().regex(/^\d{6}$/);

  type Props = {
    zid: string;
    onsuccess: (data: SubmitOtpResponseBody) => void;
    oncancel: () => void;
  };

  let { zid, onsuccess, oncancel }: Props = $props();

  let otpError = $state<string | null>(null);
  let otp = $state("");

  async function handleSubmitOTP() {
    otpError = null;
    try {
      OTP_VALIDATOR.parse(otp);
      const { data, error } = await submitOtp({ body: { zid, otp }, credentials: "include" });
      if (error) {
        otpError = error.detail ?? "Invalid OTP.";
        return;
      }
      onsuccess(data);
    } catch (e) {
      if (e instanceof z.ZodError) {
        otpError = "Invalid OTP.";
      }
    }
  }
</script>

<Card.Header>
  <Card.Title>Enter OTP</Card.Title>
  <Card.Description
    >We sent an OTP to your UNSW Email. Please enter it below to proceed.</Card.Description
  >
</Card.Header>
<Card.Content>
  <form onsubmit={handleSubmitOTP} class="flex w-full flex-col gap-1.5">
    <InputOTP.Root maxlength={6} bind:value={otp}>
      {#snippet children({ cells })}
        <InputOTP.Group>
          {#each cells as cell (cell)}
            <InputOTP.Slot {cell} aria-invalid={otpError !== null} />
          {/each}
        </InputOTP.Group>
      {/snippet}
    </InputOTP.Root>
    {#if otpError !== null}
      <p class="text-sm text-muted-foreground">{otpError}</p>
    {/if}
    <Button class="gap-1" type="submit">
      <span>Submit</span>
      <ArrowRight />
    </Button>
    <Button variant="outline" onclick={oncancel}>Use a different zID</Button>
  </form>
</Card.Content>
