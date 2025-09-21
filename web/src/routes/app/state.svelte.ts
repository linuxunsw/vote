import { state, type StateChangeEvent } from "$lib/api";

export class State {
  state = $state<StateChangeEvent["new_state"] | null>(null);

  constructor() {
    $effect(() => {
      // TODO: use liam's new state endpoint when it gets merged
      this.state = "NOMINATIONS_OPEN";

      // let generator: Awaited<ReturnType<typeof state>>["stream"];
      //
      // (async () => {
      //   const res = await state();
      //   generator = res.stream;
      //   for await (const received of generator) {
      //     // @ts-ignore TODO: fix the types for this, this code works
      //     this.state = received.new_state;
      //     console.log(this.state);
      //   }
      // })();

      // TODO: is some cleanup function required...?
    });
  }
}
