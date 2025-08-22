<script lang="ts">
  import "../app.css";

  let { children } = $props();
</script>

<div class="flex min-h-screen flex-col">
  <header class="bg-yellow-800 p-4">Vote</header>
  <div class="container mx-auto flex grow flex-col">
    <main class="grow space-y-4 p-4">
      {@render children?.()}
    </main>
  </div>
  <footer class="bg-neutral-800 p-4">&copy; Linux Society 2025</footer>
</div>
