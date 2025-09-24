<script lang="ts">
  import ElectionStateForm from "./ElectionStateForm.svelte";
  import MemberForm from "./MemberForm.svelte";
  import NewElectionForm from "./NewElectionForm.svelte";

  let { data } = $props();
</script>

<div class="space-y-6">
  <p>Current Election State: <code>{data.state.state}</code></p>

  {#if data.state.state === "NO_ELECTION"}
    <NewElectionForm />
  {:else}
    <ElectionStateForm currentState={data.state.state} />
    <!-- if state is not NO_ELECTION, then election_id is guaranteed to exist -->
    <MemberForm election_id={data.state.election_id!} />
  {/if}
</div>
