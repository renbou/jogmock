<script lang="ts">
  import type { z, ZodSchema } from "zod";
  import type { FormConfig } from "@felte/core";
  import { createForm } from "felte";
  import { validator } from "@felte/validator-zod";
  import reporter from "@felte/reporter-dom";

  // Experimental generics using
  // https://github.com/dummdidumm/rfcs/blob/ts-typedefs-within-svelte-components/text/ts-typing-props-slots-events.md
  type T = $$Generic<ZodSchema>;
  type Data = z.infer<T>;
  type Config = FormConfig<Data>;

  export let schema: T;
  export let onSubmit: Config["onSubmit"];

  const { form } = createForm<Data>({
    onSubmit,
    extend: [
      validator({ schema }),
      reporter({
        single: true,
        singleAttributes: { class: ["help", "is-danger"] },
      }),
    ],
  });
</script>

<!--
  @component Form component with validation using Felte + Zod.
-->

<form class="box" use:form>
  <slot />
</form>
