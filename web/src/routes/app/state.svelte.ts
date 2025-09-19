export type ElectionState =
  | "CLOSED"
  | "NOMINATIONS_OPEN"
  | "NOMINATIONS_CLOSED"
  | "VOTING_OPEN"
  | "VOTING_CLOSED"
  | "RESULTS"
  | "END";

export class State {
  state = $state<ElectionState | null>(null);

  constructor() {
    $effect(() => {
      const source = new EventSource("/");
      source.onmessage = () => {
        // TODO: set state here
      };

      // while developing ui, just change this to force a specific election state
      this.state = "NOMINATIONS_OPEN";

      return () => source.close();
    });
  }
}
