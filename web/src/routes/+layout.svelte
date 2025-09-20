<script lang="ts">
  import "../app.css";
  import { ModeWatcher } from "mode-watcher";
  import SEO from "$lib/components/SEO.svelte";
  import ModeToggle from "$lib/components/mode-toggle.svelte";
  import { authState } from "$lib/stores.svelte";

  let { children } = $props();

  $effect(() => {
    authState.load();
  });
</script>

<ModeWatcher />

<SEO title="LinSoc Vote" />

<div class="flex min-h-screen flex-col">
  <header class="p-4 font-bold">Vote</header>
  <div class="mx-auto w-full grow space-y-4 p-4 md:w-2/3">
    {@render children?.()}
  </div>
  <footer class="flex flex-row content-between items-center justify-center gap-4 p-4">
    <p>Linux Society 2025</p>
    <ModeToggle />
  </footer>
</div>
