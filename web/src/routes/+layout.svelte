<script lang="ts">
  import "../app.css";
  import { ModeWatcher } from "mode-watcher";
  import SEO from "$lib/components/SEO.svelte";
  import { authState } from "$lib/stores.svelte";
  import SiteFooter from "$lib/components/site-footer.svelte";
  import SiteHeader from "$lib/components/site-header.svelte";

  let { children } = $props();

  $effect(() => {
    authState.load();
  });
</script>

<ModeWatcher />

<SEO title="LinSoc Vote" />

<div class="flex min-h-svh flex-col">
  <SiteHeader />
  <main class="mx-auto flex flex-1 flex-col space-y-4 p-4 md:w-2/3">
    {@render children?.()}
  </main>
  <SiteFooter />
</div>
