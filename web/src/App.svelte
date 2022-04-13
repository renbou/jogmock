<script lang="ts">
  import router from "@lib/router";
</script>

<div class="flex flex-col h-screen">
  <header>
    <nav class="navbar is-spaced is-primary">
      <div class="navbar-brand">
        <a class="navbar-item bg-primary font-bold text-2xl" href="/">Jogmock</a
        >
      </div>
    </nav>
  </header>
  <!-- @Render currently enabled route -->
  <div class="grow"><svelte:component this={$router} /></div>
  <footer class="footer shrink-0">
    <div class="content text-center">
      <!-- prettier-ignore -->
      <p><strong>Jogmock</strong> by <a href="https://github.com/renbou">renbou</a></p>
    </div>
  </footer>
</div>
