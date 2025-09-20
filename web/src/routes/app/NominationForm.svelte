<script lang="ts">
  import * as Form from "$lib/components/ui/form/index.js";
  import { Input } from "$lib/components/ui/input";
  import { Textarea } from "$lib/components/ui/textarea";
  import { Checkbox } from "$lib/components/ui/checkbox";
  import { defaults, superForm } from "sveltekit-superforms";
  import { zod4Client, zodClient } from "sveltekit-superforms/adapters";
  import { z } from "zod";

  const execRoleItems = [
    { id: "president", label: "President" },
    { id: "secretary", label: "Secretary" },
    { id: "treasurer", label: "Treasurer" },
    { id: "arc_delegate", label: "Arc Delegate" },
    { id: "edi_officer", label: "EDI Officer" },
    { id: "grievance_officer", label: "Grievance Officer" },
  ] as const;

  const schema = z.object({
    candidate_name: z.string().min(2),
    contact_email: z.email(),
    discord_username: z.string(),
    executive_roles: z.array(z.string()),
    candidate_statement: z.string(),
    url: z.url().optional(),
  });

  type Props = {
    // TODO: make this a concrete type (probably from the codegen api)
    nomination?: {
      candidate_name: string;
      candidate_statement: string;
      contact_email: string;
      discord_username: string;
      executive_roles: string[];
      url?: string;
    };

    onsuccess: () => void;
    oncancel: () => void;
  };

  let { nomination, onsuccess, oncancel }: Props = $props();

  const form = superForm(defaults(nomination, zod4Client(schema)), {
    SPA: true,
    validators: zod4Client(schema),
    onUpdate({ form }) {
      if (!form.valid) return;
      // TODO
      onsuccess();
    },
  });
  const { form: formData, enhance } = form;

  function addExecRole(id: string) {
    $formData.executive_roles = [...$formData.executive_roles, id];
  }

  function removeExecRole(id: string) {
    $formData.executive_roles = $formData.executive_roles.filter((i) => i !== id);
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
        {@const checked = $formData.executive_roles.includes(item.id)}
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
