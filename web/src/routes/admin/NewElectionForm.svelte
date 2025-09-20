<script lang="ts">
  import { Button } from "$lib/components/ui/button";
  import { Label } from "$lib/components/ui/label";
  import { Input } from "$lib/components/ui/input";

  let name = $state("");

  function handleSubmit(e: SubmitEvent) {
    e.preventDefault();
  }
</script>

<form onsubmit={handleSubmit} class="space-y-1.5">
  <h1 class="text-2xl font-bold">Create New Election</h1>
  <Label>Election Name</Label>
  <Input bind:value={name} required />
  <Button type="submit">Create Election</Button>
</form>
