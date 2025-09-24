import { getElectionState, type GetElectionStateResponse } from "$lib/api";

export class State {
  state = $state<GetElectionStateResponse["state"] | null>(null);

  constructor() {
    $effect(() => {
      getElectionState().then(({ data, error }) => {
        if (data) this.state = data.state;
        // TODO: better error display
        if (error) throw error;
      });

      // TODO: SSE?

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
