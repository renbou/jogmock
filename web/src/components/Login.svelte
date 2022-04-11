<script lang="ts">
  import router, { Route } from "@lib/router";

  let form: HTMLFormElement;
  const info = {
    email: "",
    password: "",
    captcha: "",
  };

  let validationErrors = { ...info };
  function validate(): boolean {
    for (const key in info) {
      const input = form[key] as HTMLInputElement;
      validationErrors[key as keyof typeof validationErrors] =
        input.checkValidity() ? "" : input.validationMessage;
    }
    return false;
  }

  function submit() {
    validate();
  }
</script>

<section class="section">
  <div class="columns is-centered">
    <div class="column is-one-third">
      <form class="box" bind:this={form}>
        <div class="field">
          <label class="label" for="email">Email</label>
          <div class="control">
            <input
              class={["input", validationErrors.email && "is-danger"].join(" ")}
              name="email"
              type="email"
              required
              placeholder="email@example.com"
              bind:value={info.email}
            />
          </div>
          {#if validationErrors.email !== ""}
            <p class="help is-danger">
              {validationErrors.email}
            </p>
          {/if}
        </div>
        <div class="field">
          <label class="label" for="password">Password</label>
          <div class="control">
            <input
              class={["input", validationErrors.password && "is-danger"].join(
                " "
              )}
              name="password"
              type="password"
              required
              placeholder="••••••••••"
              bind:value={info.password}
            />
          </div>
          {#if validationErrors.password !== ""}
            <p class="help is-danger">
              {validationErrors.password}
            </p>
          {/if}
        </div>
        <div class="field">
          <label class="label" for="captcha">reCAPTCHA Token</label>
          <div class="control">
            <input
              class={["input", validationErrors.captcha && "is-danger"].join(
                " "
              )}
              name="captcha"
              type="text"
              required
              placeholder="Strava's reCAPTCHA token"
              bind:value={info.captcha}
            />
          </div>
          {#if validationErrors.captcha !== ""}
            <p class="help is-danger">
              {validationErrors.captcha}
            </p>
          {/if}
        </div>
        <div class="field">
          <div class="control">
            <input
              type="submit"
              class="button is-primary is-light is-outlined px-5"
              value="Login"
              on:click|preventDefault={submit}
            />
          </div>
        </div>
      </form>
    </div>
  </div>
</section>
