<script lang="ts">
  import { LoaderCircle } from "@lucide/svelte";

  type Props = {
    message: string;
  };

  let { message }: Props = $props();
</script>

<div class="grid h-full">
  <div class="space-y-2 place-self-center text-center">
    <div class="flex w-full justify-center">
      <LoaderCircle size={54} class="animate-spin" />
    </div>
    <h2 class="text-2xl font-bold">{message}</h2>
    <h3 class="text-lg">
      in the meantime, wanna follow our <a
        class="text-blue-500 underline"
        href="https://instagram.com/linsocunsw">instagram</a
      >?
    </h3>
  </div>
</div>
