<script lang="ts">
  import ElectionStateForm from "./ElectionStateForm.svelte";
  import MemberForm from "./MemberForm.svelte";
  import NewElectionForm from "./NewElectionForm.svelte";

  let { data } = $props();
</script>

<p>Current Election State: <code>{data.state.state}</code></p>

{#if data.state.state === "NO_ELECTION"}
  <NewElectionForm />
{:else}
  <ElectionStateForm />
  <MemberForm />
{/if}
