<script lang="ts">
  import { nanoid } from "nanoid";

  export let label: string;
  export let name: string;
  let cls: string = "";
  export { cls as class };

  const validationId = nanoid();
  const inputId = nanoid();
</script>

<!--
  @component Input field component with validation support using Felte's reporter-dom.
-->

<div class="field">
  <label class="label" for={inputId}>{label}</label>
  <div class="control">
    <input
      id={inputId}
      {name}
      class={`input ${cls}`}
      aria-describedby={validationId}
      {...$$restProps}
    />
  </div>
  <div id={validationId} data-felte-reporter-dom-for={name} />
</div>
