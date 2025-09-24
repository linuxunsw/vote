<script lang="ts">
  import { Clock } from "@lucide/svelte";
  import * as Card from "$lib/components/ui/card";

  type Props = {
    message: string;
    description: string;
  };
  let { message, description }: Props = $props();
</script>

<div class="flex h-full items-center justify-center p-4 sm:p-6 md:p-8">
  <Card.Root class="w-full max-w-md text-center">
    <Card.Content class="p-6">
      <div class="mx-auto flex h-16 w-16 items-center justify-center">
        <Clock class="h-12 w-12 text-primary" />
      </div>
      <h2 class="mt-4 text-2xl font-bold tracking-tight">{message}</h2>
      <p class="mt-2 text-muted-foreground">{description}</p>
    </Card.Content>
  </Card.Root>
</div>
