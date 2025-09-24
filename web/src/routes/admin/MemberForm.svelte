<script lang="ts">
  import { setElectionMembers, type ErrorModel } from "$lib/api";
  import { Button } from "$lib/components/ui/button";
  import { Input } from "$lib/components/ui/input";
  import { Label } from "$lib/components/ui/label";
  import { toast } from "svelte-sonner";

  type Props = {
    election_id: string;
  };

  let { election_id }: Props = $props();

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
      const matches = [...(text.match(/z\d{9}/g) ?? [])];

      toast.promise(
        setElectionMembers({ path: { election_id }, body: { zids: matches }, throwOnError: true }),
        {
          loading: "Updating member list...",
          success: "Member list successfully updated!",
          error: (e) => {
            return (e as ErrorModel).detail ?? "An error occurred.";
          },
        },
      );
    };

    reader.onerror = (e) => {
      toast.error("Error reading file: " + e.toString());
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
