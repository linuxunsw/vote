<script lang="ts">
  import "../app.css";
  import SEO from "$lib/components/SEO.svelte";

  let { children } = $props();
</script>

<SEO title="LinSoc Vote" />

<div class="flex min-h-screen flex-col">
  <header class="bg-yellow-500 p-4 font-bold">Vote</header>
  <div class="mx-auto w-full grow space-y-4 p-4 md:w-2/3">
    {@render children?.()}
  </div>
  <footer class="bg-neutral-300 p-4">&copy; Linux Society 2025</footer>
</div>
