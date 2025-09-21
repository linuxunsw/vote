<script lang="ts">
  import * as Form from "$lib/components/ui/form/index.js";
  import { Input } from "$lib/components/ui/input";
  import { Textarea } from "$lib/components/ui/textarea";
  import { Checkbox } from "$lib/components/ui/checkbox";
  import { defaults, superForm } from "sveltekit-superforms";
  import { zod4Client } from "sveltekit-superforms/adapters";
  import {
    submitNomination,
    zSubmitNominationWritable,
    type SubmitNominationWritable,
  } from "$lib/api";

  const superFormsAdapter = zod4Client(zSubmitNominationWritable);

  const execRoleItems = [
    { id: "president", label: "President" },
    { id: "secretary", label: "Secretary" },
    { id: "treasurer", label: "Treasurer" },
    { id: "arc_delegate", label: "Arc Delegate" },
    { id: "edi_officer", label: "EDI Officer" },
    { id: "grievance_officer", label: "Grievance Officer" },
  ] as const;

  type Props = {
    nomination?: SubmitNominationWritable;
    onsuccess: () => void;
    oncancel: () => void;
  };

  let {
    nomination = {
      candidate_name: "",
      contact_email: "",
      discord_username: "",
      executive_roles: [],
      candidate_statement: "",
      url: "",
    },
    onsuccess,
    oncancel,
  }: Props = $props();

  const form = superForm(defaults(nomination, superFormsAdapter), {
    SPA: true,
    validators: superFormsAdapter,
    async onUpdate({ form }) {
      if (!form.valid) return;
      const { error } = await submitNomination({ body: form.data });
      if (error) {
        // TODO
        return;
      }
      onsuccess();
    },
  });
  const { form: formData, enhance } = form;

  function addExecRole(id: NonNullable<SubmitNominationWritable["executive_roles"]>[number]) {
    $formData.executive_roles = [...($formData.executive_roles ?? []), id];
  }

  function removeExecRole(id: string) {
    $formData.executive_roles = ($formData.executive_roles ?? []).filter((i) => i !== id);
  }
</script>

<form method="POST" use:enhance>
  <Form.Field {form} name="candidate_name">
    <Form.Control>
      {#snippet children({ props })}
        <Form.Label>Name</Form.Label>
        <Input {...props} bind:value={$formData.candidate_name} />
      {/snippet}
    </Form.Control>
    <Form.FieldErrors />
  </Form.Field>

  <Form.Field {form} name="contact_email">
    <Form.Control>
      {#snippet children({ props })}
        <Form.Label>Email</Form.Label>
        <Input {...props} bind:value={$formData.contact_email} />
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

  <Form.Fieldset {form} name="executive_roles" class="space-y-0">
    <Form.Legend class="mb-1.5 text-sm">Nominating For</Form.Legend>
    <div class="space-y-2">
      {#each execRoleItems as item (item.id)}
        {@const checked = ($formData.executive_roles ?? []).includes(item.id)}
        <div class="flex flex-row items-start space-x-3">
          <Form.Control>
            {#snippet children({ props })}
              <Checkbox
                {...props}
                {checked}
                value={item.id}
                onCheckedChange={(v) => {
                  if (v) {
                    addExecRole(item.id);
                  } else {
                    removeExecRole(item.id);
                  }
                }}
              />
              <Form.Label class="font-normal">
                {item.label}
              </Form.Label>
            {/snippet}
          </Form.Control>
        </div>
      {/each}
      <Form.FieldErrors />
    </div>
  </Form.Fieldset>

  <Form.Field {form} name="candidate_statement">
    <Form.Control>
      {#snippet children({ props })}
        <Form.Label>Candidate Statement</Form.Label>
        <Textarea {...props} bind:value={$formData.candidate_statement} />
      {/snippet}
    </Form.Control>
    <Form.FieldErrors />
  </Form.Field>

  <Form.Field {form} name="url">
    <Form.Control>
      {#snippet children({ props })}
        <Form.Label>URL (optional)</Form.Label>
        <Input {...props} bind:value={$formData.url} />
      {/snippet}
    </Form.Control>
    <Form.FieldErrors />
  </Form.Field>
  <Form.Button>Submit</Form.Button>
  <Form.Button variant="outline" type="button" onclick={oncancel}>Cancel</Form.Button>
</form>
