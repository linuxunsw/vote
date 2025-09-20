<script lang="ts">
  import ClosedState from "./ClosedState.svelte";
  import InterimState from "./InterimState.svelte";
  import NominationState from "./NominationState.svelte";
  import { State } from "./state.svelte";

  let { data } = $props();

  let state = new State();
</script>

{#if state.state === "CLOSED"}
  <ClosedState
    message="Looks like no elections are running right now..."
    imageSrc="https://i.pinimg.com/736x/d9/1d/aa/d91daa69daf44b52728271bd478898d3.jpg"
  />
{:else if state.state === "NOMINATIONS_OPEN"}
  <NominationState nomination={data.nomination} />
{:else if state.state === "NOMINATIONS_CLOSED"}
  <InterimState message="Nominations Closed! Voting will open shortly..." />
{:else if state.state === "VOTING_OPEN"}
  <!-- TODO -->
{:else if state.state === "VOTING_CLOSED"}
  <InterimState message="Waiting on results... Hang tight!" />
{:else if state.state === "RESULTS"}
  <!-- TODO -->
{:else if state.state === "END"}
  <!-- TODO -->
{/if}
