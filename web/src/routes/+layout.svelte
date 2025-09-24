<script lang="ts">
  import "../app.css";
  import { ModeWatcher } from "mode-watcher";
  import { Toaster } from "$lib/components/ui/sonner/index.js";
  import SEO from "$lib/components/SEO.svelte";
  import { authState } from "$lib/stores.svelte";
  import SiteFooter from "$lib/components/site-footer.svelte";
  import SiteHeader from "$lib/components/site-header.svelte";

  let { children } = $props();

  $effect(() => {
    authState.load();
  });
</script>

<Toaster richColors />
<ModeWatcher />

<SEO title="LinSoc Vote" />

<div class="flex min-h-svh flex-col">
  <SiteHeader />
  <main
    class="p-y-0 p-x-4 sm:p-x-6 md:p-x-8 mx-auto flex max-w-[calc(var(--breakpoint-2xl)+2rem)] flex-1 flex-col"
  >
    {@render children?.()}
  </main>
  <SiteFooter />
</div>
