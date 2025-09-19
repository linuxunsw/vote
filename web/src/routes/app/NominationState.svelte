<script lang="ts">
  import Button from "$lib/components/ui/button/button.svelte";
  import NominationForm from "./NominationForm.svelte";

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
    // TODO: make this a concrete type (probably from the codegen api)
    nomination?: {
      candidate_name: string;
      candidate_statement: string;
      contact_email: string;
      discord_username: string;
      executive_roles: string[];
      url?: string;
    };
  };

  let { nomination }: Props = $props();
  let editing = $state(false);
</script>

{#if nomination}
  <p>You are currently nominated in this election.</p>
{:else}
  <p>You are not currently nominated in this election.</p>
{/if}

{#if !editing}
  {#if nomination}
    <div class="grid grid-cols-1 gap-x-1.5 gap-y-1 md:grid-cols-[auto_1fr] md:gap-y-2">
      <p class="font-bold md:text-right">Name</p>
      <p class="mb-1.5 md:mb-0">{nomination.candidate_name}</p>

      <p class="font-bold md:text-right">Contact Email</p>
      <a class="mb-1.5 text-blue-600 underline md:mb-0" href="mailto:{nomination.contact_email}">
        {nomination.contact_email}
      </a>

      <p class="font-bold md:text-right">Discord Username</p>
      <p class="mb-1.5 md:mb-0">{nomination.discord_username}</p>

      <p class="font-bold md:text-right">Nominated For</p>
      <p class="mb-1.5 md:mb-0">{nomination.executive_roles.map(humanizeExecRole).join(", ")}</p>

      <p class="font-bold md:text-right">Candidate Statement</p>
      <p class="mb-1.5 md:mb-0">{nomination.candidate_statement}</p>

      <p class="font-bold md:text-right">URL</p>
      {#if nomination.url}
        <!-- eslint-disable-next-line svelte/no-navigation-without-resolve -->
        <a class="mb-1.5 text-blue-600 underline md:mb-0" href={nomination.url}>{nomination.url}</a>
      {:else}
        <p class="mb-1.5 text-muted-foreground italic md:mb-0">(not provided)</p>
      {/if}
    </div>

    <div class="flex gap-2">
      <Button onclick={() => (editing = true)}>Edit Nomination</Button>
      <Button variant="destructive">Revoke Nomination</Button>
    </div>
  {:else}
    <Button onclick={() => (editing = true)}>Nominate Yourself!</Button>
  {/if}
{:else}
  <NominationForm {nomination} onsuccess={() => {}} oncancel={() => (editing = false)} />
{/if}
