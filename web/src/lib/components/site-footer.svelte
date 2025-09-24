<script lang="ts">
  import { authState } from "$lib/stores.svelte";
</script>

<footer class="bg-transparent">
  <div class="mx-auto max-w-4xl px-6">
    <div class="flex h-24 items-center justify-center gap-2">
      <div class="text-center text-xs leading-loose text-muted-foreground sm:text-sm">
        Built by the
        <a
          href="https://linuxunsw.org"
          target="_blank"
          rel="noreferrer"
          class="font-medium underline underline-offset-4"
        >
          Linux Society @ UNSW
        </a>
      </div>
      {#if authState.value?.is_admin}
        <span class="text-muted-foreground">|</span>
        <div class="text-center text-xs leading-loose text-muted-foreground sm:text-sm">
          <a href="/admin" class="font-medium underline underline-offset-4">Admin Dashboard</a>
        </div>
      {/if}
    </div>
  </div>
</footer>
