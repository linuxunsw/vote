import type { SubmitOtpResponseBody } from "./api";

class PersistedState<T> {
  #key: string;
  #value = $state<T>();

  constructor(key: string) {
    this.#key = key;
  }

  load(initialData?: T) {
    const saved = localStorage.getItem(this.#key);
    if (saved) {
      this.#value = JSON.parse(saved);
    } else {
      this.value = initialData;
    }
  }

  set value(v: T | undefined) {
    this.#value = v;
    if (v) {
      localStorage.setItem(this.#key, JSON.stringify(v));
    } else {
      localStorage.removeItem(this.#key);
    }
  }

  get value() {
    return this.#value;
  }
}

export const authState = new PersistedState<SubmitOtpResponseBody>("auth");
