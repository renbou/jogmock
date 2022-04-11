<script lang="ts">
  import router, { Route } from "@lib/router";
</script>

<h1>home</h1>
<button on:click={() => router.navigate(Route.Login)}>login</button>
