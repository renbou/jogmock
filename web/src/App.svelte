<script lang="ts">
  import router from '@lib/router';
</script>

<div class="is-flex is-flex-direction-column height-100vh">
  <header>
    <nav class="navbar is-spaced is-primary">
      <div class="navbar-brand">
        <a
          class="navbar-item has-background-primary has-text-weight-bold is-size-4 brand"
          href="/">Jogmock</a
        >
      </div>
    </nav>
  </header>
  <!-- @Render currently enabled route -->
  <div class="is-flex-grow-1"><svelte:component this={$router} /></div>
  <footer class="footer is-flex-shrink-0">
    <div class="content has-text-centered">
      <!-- prettier-ignore -->
      <p><strong>Jogmock</strong> by <a href="https://github.com/renbou">renbou</a></p>
    </div>
  </footer>
</div>

<style lang="scss">
  .brand {
    border-radius: $radius-large;
  }

  .height-100vh {
    height: 100vh;
  }
</style>
