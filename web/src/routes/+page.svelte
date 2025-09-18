<script lang="ts">
  import Separator from "$lib/components/ui/separator/separator.svelte";
  import About from "./About.svelte";
  import LoginCard from "./LoginCard.svelte";
</script>

<div class="mx-auto grid h-full md:max-w-2/3">
  <div
    class="grid grid-rows-[1fr_auto_1fr] gap-5 place-self-center md:grid-cols-[1fr_auto_1fr] md:grid-rows-1"
  >
    <LoginCard />
    <Separator class="md:hidden" />
    <Separator orientation="vertical" class="hidden md:inline" />
    <About />
  </div>
</div>
