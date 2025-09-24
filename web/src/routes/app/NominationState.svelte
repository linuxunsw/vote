<script lang="ts">
  import { deleteNomination, type Nomination } from "$lib/api";
  import { invalidateAll } from "$app/navigation";
  import { Button, buttonVariants } from "$lib/components/ui/button";
  import * as Card from "$lib/components/ui/card";
  import * as Popover from "$lib/components/ui/popover";
  import NominationForm from "./NominationForm.svelte";
  import { SquarePen, Trash2 } from "@lucide/svelte";

  function humaniseExecRole(role: string) {
    return (
      {
        president: "President",
        secretary: "Secretary",
        treasurer: "Treasurer",
        arc_delegate: "Arc Delegate",
        grievance_officer: "Grievance Officer",
        edi_officer: "EDI Officer",
      }[role] ?? role
    );
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

{#if editing}
  <NominationForm {nomination} onsuccess={handleSubmit} oncancel={() => (editing = false)} />
{:else if nomination}
  <Card.Root class="max-w-128">
    <Card.Header class="grid grid-cols-1 items-start gap-4 sm:grid-cols-2">
      <div class="space-y-2">
        <Card.Title class="text-2xl">Your Nomination</Card.Title>
        <Card.Description>You are currently nominated in this election.</Card.Description>
      </div>

      <div class="flex items-start justify-start sm:items-end sm:justify-end">
        <Button variant="outline" onclick={() => (editing = true)}>
          <SquarePen class="h-4 w-4" />
          <span>Edit Nomination</span>
        </Button>
      </div>
    </Card.Header>
    <Card.Content class="grid gap-6">
      <dl class="grid grid-cols-1 gap-x-4 gap-y-6 sm:grid-cols-2">
        <div class="sm:col-span-1">
          <dt class="text-sm font-medium text-muted-foreground">Name</dt>
          <dd class="mt-1 text-lg font-semibold">{nomination.candidate_name}</dd>
        </div>
        <div class="sm:col-span-1">
          <dt class="text-sm font-medium text-muted-foreground">Nominated For</dt>
          <dd class="mt-1 text-lg font-semibold">
            {(nomination.executive_roles ?? []).map(humaniseExecRole).join(", ")}
          </dd>
        </div>
        <div class="sm:col-span-1">
          <dt class="text-sm font-medium text-muted-foreground">Contact Email</dt>
          <dd class="mt-1">
            <a class="underline-offset-4 hover:underline" href="mailto:{nomination.contact_email}">
              {nomination.contact_email}
            </a>
          </dd>
        </div>
        <div class="sm:col-span-1">
          <dt class="text-sm font-medium text-muted-foreground">Discord Username</dt>
          <dd class="mt-1">{nomination.discord_username}</dd>
        </div>
        <div class="sm:col-span-2">
          <dt class="text-sm font-medium text-muted-foreground">Candidate Statement</dt>
          <dd class="mt-1 break-words whitespace-pre-wrap">{nomination.candidate_statement}</dd>
        </div>
        <div class="sm:col-span-2">
          <dt class="text-sm font-medium text-muted-foreground">URL</dt>
          <dd class="mt-1">
            {#if nomination.url}
              <!-- eslint-disable svelte/no-navigation-without-resolve -->
              <a
                class="underline-offset-4 hover:underline"
                href={nomination.url}
                target="_blank"
                rel="noreferrer">{nomination.url}</a
              >
              <!-- eslint-enable svelte/no-navigation-without-resolve -->
            {:else}
              <span class="text-muted-foreground italic">(not provided)</span>
            {/if}
          </dd>
        </div>
      </dl>
    </Card.Content>
    <Card.Footer class="justify-end">
      <Popover.Root>
        <Popover.Trigger class={buttonVariants({ variant: "destructive" })}>
          <Trash2 class="mr-2 h-4 w-4" />
          Revoke Nomination
        </Popover.Trigger>
        <Popover.Content class="w-fit space-y-3 p-4">
          <p class="font-semibold">Are you sure?</p>
          <Button onclick={handleRevoke} variant="destructive" class="w-full">Confirm</Button>
        </Popover.Content>
      </Popover.Root>
    </Card.Footer>
  </Card.Root>
{:else}
  <Card.Root class="max-w-128">
    <Card.Header>
      <Card.Title class="text-2xl">Get in the running!</Card.Title>
      <Card.Description>Nominations are open, but you haven't submitted one yet.</Card.Description>
    </Card.Header>
    <Card.Content>
      <p>Submit your nomination to run for an executive role in the upcoming election.</p>
    </Card.Content>
    <Card.Footer>
      <Button onclick={() => (editing = true)}>Nominate Yourself</Button>
    </Card.Footer>
  </Card.Root>
{/if}
