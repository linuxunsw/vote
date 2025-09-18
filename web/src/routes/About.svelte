<script lang="ts">
  import { Github, Instagram, BotMessageSquare, University } from "@lucide/svelte";
  import type { Component } from "svelte";
</script>

{#snippet social(label: string, href: string, Icon: Component)}
  <a {href} class="flex gap-1.5 text-muted-foreground underline hover:text-foreground">
    <Icon />
    <span>{label}</span>
  </a>
{/snippet}

<div class="flex flex-col gap-3">
  <p class="text-lg">
    <strong>Vote</strong> is a web application designed to streamline the management of elections in
    UNSW societies, starting with Linux Society.
  </p>
  <div class="flex flex-col gap-2">
    {@render social("Source", "https://github.com/linuxunsw/vote", Github)}
    {@render social(
      "spARC",
      "https://member.arc.unsw.edu.au/s/clubdetail?clubid=0016F00004CaG1EQAV",
      University,
    )}
    {@render social("Discord", "https://discord.gg/7Tqz2bn8Gz", BotMessageSquare)}
    {@render social("Instagram", "https://instagram.com/linsocunsw", Instagram)}
  </div>
</div>
