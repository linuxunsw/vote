<script lang="ts">
  type Props = {
    message: string;
    imageSrc: string;
  };

  let { message, imageSrc }: Props = $props();
</script>

<div class="grid h-full">
  <div class="space-y-2 place-self-center text-center">
    <img alt="funny penguin" class="mx-auto max-w-1/2" src={imageSrc} />
    <h2 class="text-2xl font-bold">{message}</h2>
    <h3 class="text-lg">
      in the meantime, wanna follow our <a
        class="text-blue-500 underline"
        href="https://instagram.com/linsocunsw">instagram</a
      >?
    </h3>
  </div>
</div>
