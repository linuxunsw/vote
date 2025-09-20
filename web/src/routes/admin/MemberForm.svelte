<script lang="ts">
  import { Button } from "$lib/components/ui/button";
  import { Input } from "$lib/components/ui/input";
  import { Label } from "$lib/components/ui/label";

  let files = $state<FileList>();

  function handleSubmit(e: SubmitEvent) {
    e.preventDefault();

    if (!files) return;

    let file = files[0];

    const reader = new FileReader();
    reader.onload = (e) => {
      if (!e.target) return;
      const text = e.target.result;
      if (typeof text !== "string") return;
      const matches = text.match(/z\d{9}/g);

      // TODO: api stuff
      console.log(matches);

      // TODO: success feedback
    };

    reader.onerror = (e) => {
      // TODO: error feedback
    };

    reader.readAsText(file);
  }
</script>

<form onsubmit={handleSubmit} class="space-y-1.5">
  <h1 class="text-2xl font-bold">Manage Member List</h1>
  <Label>Upload zID File</Label>
  <Input type="file" bind:files required />
  <Button type="submit">Update Member List</Button>
</form>
