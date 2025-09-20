<script lang="ts">
  import { Button } from "$lib/components/ui/button";
  import { Label } from "$lib/components/ui/label";
  import * as Select from "$lib/components/ui/select";

  const POSSIBLE_STATES = [
    "CLOSED",
    "NOMINATIONS_OPEN",
    "NOMINATIONS_CLOSED",
    "VOTING_OPEN",
    "VOTING_CLOSED",
    "RESULTS",
    "END",
  ] as const;

  let newState = $state("CLOSED");

  function handleSubmit(e: SubmitEvent) {
    e.preventDefault();
  }
</script>

<form onsubmit={handleSubmit} class="space-y-1.5">
  <h1 class="text-2xl font-bold">Manage Election State</h1>
  <Label>New Election State</Label>
  <Select.Root type="single" bind:value={newState}>
    <Select.Trigger>{newState}</Select.Trigger>
    <Select.Content>
      {#each POSSIBLE_STATES as state (state)}
        <Select.Item value={state} label={state}>{state}</Select.Item>
      {/each}
    </Select.Content>
  </Select.Root>
  <Button type="submit">Update State</Button>
</form>
