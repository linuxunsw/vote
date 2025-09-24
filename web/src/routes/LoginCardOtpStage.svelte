<script lang="ts">
  import { submitOtp, type SubmitOtpResponseBody } from "$lib/api";
  import Button from "$lib/components/ui/button/button.svelte";
  import * as Card from "$lib/components/ui/card";
  import * as InputOTP from "$lib/components/ui/input-otp";
  import Label from "$lib/components/ui/label/label.svelte";
  import { LoaderCircle } from "@lucide/svelte";
  import * as z from "zod";
  import { REGEXP_ONLY_DIGITS } from "bits-ui";

  const OTP_VALIDATOR = z.string().regex(/^\d{6}$/, "Your OTP must be 6 digits.");

  type Props = {
    zid: string;
    onsuccess: (data: SubmitOtpResponseBody) => void;
    oncancel: () => void;
  };

  let { zid, onsuccess, oncancel }: Props = $props();

  let otpError = $state<string | null>(null);
  let otp = $state("");
  let isLoading = $state(false);

  async function handleSubmitOTP() {
    otpError = null;
    isLoading = true;
    try {
      OTP_VALIDATOR.parse(otp);
      const { data, error } = await submitOtp({ body: { zid, otp } });
      if (error) {
        otpError = error.detail ?? "The OTP you entered is invalid or has expired.";
        return;
      }
      if (data) {
        onsuccess(data);
      }
    } catch (e) {
      if (e instanceof z.ZodError) {
        otpError = e.message;
      }
    } finally {
      isLoading = false;
    }
  }
</script>

<Card.Header class="text-center">
  <Card.Title class="text-2xl">Check your email</Card.Title>
  <Card.Description>We've sent a 6-digit code to your UNSW student email.</Card.Description>
</Card.Header>
<Card.Content>
  <form onsubmit={handleSubmitOTP} class="grid gap-4">
    <div class="grid place-items-center gap-2 text-center">
      <Label for="otp-input" class="sr-only">One-Time Password</Label>
      <InputOTP.Root id="otp-input" pattern={REGEXP_ONLY_DIGITS} maxlength={6} bind:value={otp}>
        {#snippet children({ cells })}
          <InputOTP.Group>
            {#each cells.slice(0, 3) as cell (cell)}
              <InputOTP.Slot {cell} />
            {/each}
          </InputOTP.Group>
          <InputOTP.Separator>|</InputOTP.Separator>
          <InputOTP.Group>
            {#each cells.slice(3, 6) as cell (cell)}
              <InputOTP.Slot {cell} />
            {/each}
          </InputOTP.Group>
        {/snippet}
      </InputOTP.Root>
      {#if otpError}
        <p class="text-sm text-destructive">{otpError}</p>
      {/if}
    </div>
    <div class="grid gap-2">
      <Button type="submit" class="w-full" disabled={isLoading}>
        {#if isLoading}
          <LoaderCircle class="mr-2 h-4 w-4 animate-spin" />
          <span>Verifying...</span>
        {:else}
          <span>Continue</span>
        {/if}
      </Button>
      <Button variant="link" size="sm" class="w-full" onclick={oncancel} disabled={isLoading}>
        Use a different zID
      </Button>
    </div>
  </form>
</Card.Content>
