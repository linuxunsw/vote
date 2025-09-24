<script lang="ts">
  import {
    submitNomination,
    zSubmitNominationWritable,
    type ErrorModel,
    type SubmitNominationWritable,
  } from "$lib/api";
  import * as Card from "$lib/components/ui/card";
  import { Checkbox } from "$lib/components/ui/checkbox";
  import * as Form from "$lib/components/ui/form/index.js";
  import { Input } from "$lib/components/ui/input";
  import { Textarea } from "$lib/components/ui/textarea";
  import { defaults, superForm } from "sveltekit-superforms";
  import { zodClient } from "sveltekit-superforms/adapters";
  import { toast } from "svelte-sonner";

  const superFormsAdapter = zodClient(zSubmitNominationWritable);

  const execRoleItems = [
    { id: "president", label: "President" },
    { id: "secretary", label: "Secretary" },
    { id: "treasurer", label: "Treasurer" },
    { id: "arc_delegate", label: "Arc Delegate" },
    { id: "edi_officer", label: "EDI Officer" },
    { id: "grievance_officer", label: "Grievance Officer" },
  ] as const;

  function buildSubmitPayload(data: Partial<SubmitNominationWritable>): SubmitNominationWritable {
    return {
      candidate_name: data.candidate_name ?? "",
      contact_email: data.contact_email ?? "",
      discord_username: data.discord_username ?? "",
      executive_roles: data.executive_roles ?? [],
      candidate_statement: data.candidate_statement ?? "",
      url: data.url ?? undefined,
    };
  }

  type Props = {
    nomination?: SubmitNominationWritable;
    onsuccess: () => void;
    oncancel: () => void;
  };
  let { nomination, onsuccess, oncancel }: Props = $props();

  const form = superForm(
    defaults(
      nomination ?? {
        candidate_name: "",
        contact_email: "",
        discord_username: "",
        executive_roles: [],
        candidate_statement: "",
        url: undefined,
      },
      superFormsAdapter,
    ),
    {
      validators: superFormsAdapter,
      SPA: true,
      async onUpdate({ form }) {
        if (!form.valid) return;

        const payload = buildSubmitPayload(form.data as Partial<SubmitNominationWritable>);
        toast.promise(submitNomination({ body: payload, throwOnError: true }), {
          loading: "Submitting...",
          success: () => {
            onsuccess();
            return "Successfully submitted nomination!";
          },
          error: (e) => {
            return (e as ErrorModel).detail ?? "An error has occured";
          },
        });
      },
    },
  );
  const { form: formData, enhance } = form;

  function handleCheckedChange(
    checked: boolean,
    id: NonNullable<SubmitNominationWritable["executive_roles"]>[number],
  ) {
    const currentRoles = $formData.executive_roles ?? [];
    if (checked) {
      $formData.executive_roles = [...currentRoles, id];
    } else {
      $formData.executive_roles = currentRoles.filter((role) => role !== id);
    }
  }
</script>

<Card.Root class="max-w-128">
  <Card.Header>
    <Card.Title class="text-2xl">{nomination ? "Edit Your" : "Submit a"} Nomination</Card.Title>
    <Card.Description>
      Fill out the form below to run for one or more executive roles.
    </Card.Description>
  </Card.Header>
  <Card.Content>
    <form method="POST" use:enhance class="grid gap-6">
      <div class="grid grid-cols-1 gap-6 sm:grid-cols-2">
        <Form.Field {form} name="candidate_name">
          <Form.Control>
            {#snippet children({ props })}
              <Form.Label>Full Name</Form.Label>
              <Input {...props} bind:value={$formData.candidate_name} />
            {/snippet}
          </Form.Control>
          <Form.FieldErrors />
        </Form.Field>

        <Form.Field {form} name="contact_email">
          <Form.Control>
            {#snippet children({ props })}
              <Form.Label>Contact Email</Form.Label>
              <Input {...props} type="email" bind:value={$formData.contact_email} />
            {/snippet}
          </Form.Control>
          <Form.FieldErrors />
        </Form.Field>

        <Form.Field {form} name="discord_username">
          <Form.Control>
            {#snippet children({ props })}
              <Form.Label>Discord Username</Form.Label>
              <Input {...props} bind:value={$formData.discord_username} />
            {/snippet}
          </Form.Control>
          <Form.FieldErrors />
        </Form.Field>

        <Form.Field {form} name="url">
          <Form.Control>
            {#snippet children({ props })}
              <Form.Label>URL (optional)</Form.Label>
              <Input
                {...props}
                placeholder="e.g., your personal website"
                bind:value={() => $formData.url, (v) => ($formData.url = v || undefined)}
              />
            {/snippet}
          </Form.Control>
          <Form.FieldErrors />
        </Form.Field>
      </div>

      <Form.Fieldset {form} name="executive_roles">
        <Form.Legend class="text-base font-semibold">Nominating For</Form.Legend>
        <div class="mt-2 grid grid-cols-2 gap-4 sm:grid-cols-3">
          {#each execRoleItems as item (item.id)}
            {@const checked = ($formData.executive_roles ?? []).includes(item.id)}
            <div class="flex items-center gap-x-3">
              <Checkbox
                id={`role-${item.id}`}
                {checked}
                onCheckedChange={(v) => handleCheckedChange(!!v, item.id)}
              />
              <label for={`role-${item.id}`} class="text-sm leading-none font-medium">
                {item.label}
              </label>
            </div>
          {/each}
        </div>
        <Form.FieldErrors class="mt-2" />
      </Form.Fieldset>

      <Form.Field {form} name="candidate_statement">
        <Form.Control>
          {#snippet children({ props })}
            <Form.Label>Candidate Statement</Form.Label>
            <Textarea
              {...props}
              class="min-h-[120px]"
              placeholder="Tell everyone why you'd be a great fit..."
              bind:value={$formData.candidate_statement}
            />
            <Form.Description
              class={$formData.candidate_statement.length >
              (zSubmitNominationWritable.shape.candidate_statement.maxLength ?? Infinity)
                ? "text-destructive"
                : ""}
              >{$formData.candidate_statement.length} / {zSubmitNominationWritable.shape
                .candidate_statement.maxLength} characters</Form.Description
            >
          {/snippet}
        </Form.Control>
        <Form.FieldErrors />
      </Form.Field>

      <Card.Footer class="flex justify-end gap-2 p-0 pt-4">
        <Form.Button variant="outline" type="button" onclick={oncancel}>Cancel</Form.Button>
        <Form.Button type="submit">Save Nomination</Form.Button>
      </Card.Footer>
    </form>
  </Card.Content>
</Card.Root>
