<script lang="ts">
  import { Button, buttonVariants } from "$lib/components/ui/button";
  import NominationForm from "./NominationForm.svelte";
  import * as Popover from "$lib/components/ui/popover";
  import { deleteNomination, type Nomination } from "$lib/api";
  import { invalidateAll } from "$app/navigation";

  function humanizeExecRole(role: string) {
    return {
      president: "President",
      secretary: "Secretary",
      treasurer: "Treasurer",
      arc_delegate: "Arc Delegate",
      grievance_officer: "Grievance Officer",
    }[role];
  }

  type Props = {
    nomination?: Nomination;
  };

  let { nomination }: Props = $props();
  let editing = $state(false);

  async function handleSubmit() {
    await invalidateAll();
    editing = false;
  }

  async function handleRevoke() {
    await deleteNomination();
    await invalidateAll();
  }
</script>

{#if nomination}
  <p>You are currently nominated in this election.</p>
{:else}
  <p>You are not currently nominated in this election.</p>
{/if}

{#if !editing}
  {#if nomination}
    <div class="md:space-y-3">
      <div class="space-y-1.5">
        <p class="font-bold">Name</p>
        <p>{nomination.candidate_name}</p>
      </div>

      <div class="space-y-1.5">
        <p class="font-bold">Contact Email</p>
        <a class="text-blue-600 underline" href="mailto:{nomination.contact_email}">
          {nomination.contact_email}
        </a>
      </div>

      <div class="space-y-1.5">
        <p class="font-bold">Discord Username</p>
        <p>{nomination.discord_username}</p>
      </div>

      <div class="space-y-1.5">
        <p class="font-bold">Nominated For</p>
        <p>
          {(nomination.executive_roles ?? []).map(humanizeExecRole).join(", ")}
        </p>
      </div>

      <div class="space-y-1.5">
        <p class="font-bold">Candidate Statement</p>
        <p>{nomination.candidate_statement}</p>
      </div>

      <div class="space-y-1.5">
        <p class="font-bold">URL</p>
        {#if nomination.url}
          <!-- eslint-disable-next-line svelte/no-navigation-without-resolve -->
          <a class="text-blue-600 underline" href={nomination.url}>{nomination.url}</a>
        {:else}
          <p class="text-muted-foreground italic">(not provided)</p>
        {/if}
      </div>
    </div>

    <div class="flex gap-2">
      <Button onclick={() => (editing = true)}>Edit Nomination</Button>
      <Popover.Root>
        <Popover.Trigger class={buttonVariants({ variant: "outline" })}>
          Revoke Nomination
        </Popover.Trigger>
        <Popover.Content class="w-fit space-y-1.5">
          <p class="font-bold">Are you sure?</p>
          <Button onclick={handleRevoke} variant="destructive">Confirm Revocation</Button>
        </Popover.Content>
      </Popover.Root>
    </div>
  {:else}
    <Button onclick={() => (editing = true)}>Nominate Yourself!</Button>
  {/if}
{:else}
  <!-- TODO: the api types... why... -->
  <NominationForm {nomination} onsuccess={handleSubmit} oncancel={() => (editing = false)} />
{/if}
