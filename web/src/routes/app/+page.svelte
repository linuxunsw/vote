<script lang="ts">
  import ClosedState from "./ClosedState.svelte";
  import InterimState from "./InterimState.svelte";
  import NominationState from "./NominationState.svelte";
  import { State } from "./state.svelte";

  let { data } = $props();

  let state = new State();
</script>

{#if !state.state}
  <InterimState message="Checking election status..." description="" />
{:else if state.state === "CLOSED" || state.state === "NO_ELECTION"}
  <ClosedState
    message="No elections are running right now"
    description="Check back later or follow our socials for updates."
  />
{:else if state.state === "NOMINATIONS_OPEN"}
  <NominationState nomination={data.nomination} />
{:else if state.state === "NOMINATIONS_CLOSED"}
  <InterimState
    message="Nominations Closed!"
    description="Voting will open shortly. Good luck to all candidates!"
  />
{:else if state.state === "VOTING_OPEN"}{:else if state.state === "VOTING_CLOSED"}
  <InterimState
    message="Voting is now closed!"
    description="We're counting the votes. Results will be announced soon!"
  />
{:else if state.state === "RESULTS"}{:else if state.state === "END"}{/if}
