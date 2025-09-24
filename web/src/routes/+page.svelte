<script lang="ts">
  import { toast } from "svelte-sonner";
  import About from "./About.svelte";
  import LoginCard from "./LoginCard.svelte";
  import { page } from "$app/state";

  $effect(() => {
    if (page.url.searchParams.get("error") === "session_expired") {
      toast.warning("Your session has expired, please log in again.");
    }
  });
</script>

<div class="grid h-full flex-1 md:grid-cols-2">
  <div
    class="relative hidden flex-col items-center justify-center p-10 text-muted-foreground md:flex"
  >
    <div class="absolute inset-0 bg-gradient-to-t from-background to-background/50"></div>
    <div class="relative z-20 flex flex-col items-center text-center">
      <About />
    </div>
  </div>

  <div class="flex items-center justify-center p-4">
    <div class="mx-auto flex w-full flex-col justify-center space-y-6 sm:w-[400px]">
      <LoginCard />
    </div>
  </div>
</div>
