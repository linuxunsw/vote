<script lang="ts">
  import ModeToggle from "./mode-toggle.svelte";
</script>

<header class="w-full bg-background">
  <div
    class="mx-auto flex w-full max-w-[calc(var(--breakpoint-2xl)+2rem)] items-center justify-between px-6 py-4"
  >
    <h1 class="text-3xl font-bold">Vote</h1>
    <ModeToggle />
  </div>
</header>
