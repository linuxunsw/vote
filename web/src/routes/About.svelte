<script lang="ts">
  import { simpleIconToIconNode } from "$lib/simpleIcontoIconNode";
  import { Github, Instagram, University, type IconNode, Icon } from "@lucide/svelte";
  import { siDiscord } from "simple-icons";
  import type { Component } from "svelte";
  const discordNode = simpleIconToIconNode(siDiscord, {
    title: siDiscord.title,
    fill: null,
    strokeWidth: 1.5,
  });
</script>

{#snippet social(label: string, href: string, LucideIcon?: Component, simpleIcon?: IconNode)}
  <a
    {href}
    target="_blank"
    rel="noreferrer"
    class="flex items-center gap-2 text-foreground underline-offset-4 hover:underline"
  >
    {#if LucideIcon}
      <LucideIcon class="h-5 w-5" />
    {:else if simpleIcon}
      <Icon iconNode={simpleIcon} />
    {/if}
    <span>{label}</span>
  </a>
{/snippet}

<div class="b-2 flex flex-col gap-6">
  <div class="flex flex-col gap-2">
    <h1 class="text-4xl font-bold tracking-tight text-foreground">UNSW Linux Society Elections</h1>
    <p class="text-lg text-muted-foreground">
      A secure, streamlined, and transparent voting platform.
    </p>
  </div>
  <p class="text-lg text-balance text-foreground">
    <strong>Vote</strong> is a web application designed to manage elections for UNSW societies, starting
    with LINSOC.
  </p>
  <div class="flex flex-wrap justify-center gap-x-6 gap-y-3">
    {@render social("Source", "https://github.com/linuxunsw/vote", Github)}
    {@render social("Rubric", "https://campus.hellorubric.com/?s=12591", University)}
    {@render social("Discord", "https://discord.gg/7Tqz2bn8Gz", undefined, discordNode)}
    {@render social("Instagram", "https://instagram.com/linsocunsw", Instagram)}
  </div>
</div>
