<script lang="ts">
  type Props = {
    title: string;
  };

  let { title }: Props = $props();
</script>

<svelte:head>
  <title>{title}</title>
  <meta
    name="description"
    content="Vote is an application that manages and holds elections for UNSW Linux Society."
  />
  <meta name="og:title" content={title} />
</svelte:head>
