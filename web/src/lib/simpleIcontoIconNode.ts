import type { SimpleIcon } from "simple-icons";
import type { Snippet } from "svelte";
import type { Attrs, IconNode } from "@lucide/svelte";

function makeSnippet(text: string): Snippet {
  const fn = (() => text) as unknown as Snippet;
  return fn;
}

export function simpleIconToIconNode(
  si: SimpleIcon,
  opts?: { title?: string; fill?: string | null; pad?: number; strokeWidth?: number },
): IconNode {
  const titleText = opts?.title ?? si.title;
  const wantFill = opts?.fill !== null;
  const fill = wantFill ? (opts?.fill ?? `#${si.hex}`) : "none";
  const pad = opts?.pad ?? 1; // px in 24-unit coordinates
  const strokeWidth = opts?.strokeWidth ?? 2;

  const W = 24;
  const inset = pad / W;
  const scale = 1 - inset * 2;
  const tx = W * inset;

  const nodes: IconNode = [];

  if (titleText) {
    nodes.push(["title", { children: makeSnippet(titleText) }]);
  }

  const pathAttrs: Attrs = {
    d: si.path,
    fill,
    transform: `translate(${tx} ${tx}) scale(${scale})`,
  };

  pathAttrs["stroke-width"] = strokeWidth;
  pathAttrs["stroke-linecap"] = "round";
  pathAttrs["stroke-linejoin"] = "round";

  nodes.push(["path", pathAttrs]);

  return nodes;
}
