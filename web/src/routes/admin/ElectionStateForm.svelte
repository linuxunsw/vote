<script lang="ts">
  import { invalidateAll } from "$app/navigation";
  import { adminTransitionElectionState, type ErrorModel } from "$lib/api";
  import { Button } from "$lib/components/ui/button";
  import { Label } from "$lib/components/ui/label";
  import * as Select from "$lib/components/ui/select";
  import { toast } from "svelte-sonner";

  const POSSIBLE_STATES = [
    "CLOSED",
    "NOMINATIONS_OPEN",
    "NOMINATIONS_CLOSED",
    "VOTING_OPEN",
    "VOTING_CLOSED",
    "RESULTS",
    "END",
  ] as const;

  type Props = {
    currentState: (typeof POSSIBLE_STATES)[number];
  };

  let { currentState }: Props = $props();

  let newState = $state<(typeof POSSIBLE_STATES)[number]>(currentState);

  function handleSubmit(e: SubmitEvent) {
    toast.promise(adminTransitionElectionState({ body: { state: newState }, throwOnError: true }), {
      loading: "Updating State...",
      success: () => {
        invalidateAll();
        return "State updated successfully.";
      },
      error: (e) => {
        return (e as ErrorModel).detail ?? "An error occurred.";
      },
    });
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
