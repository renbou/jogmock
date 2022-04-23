<script lang="ts">
  import router, { Route } from "@lib/router";
  import Input from "@components/Input.svelte";
  import Form from "@components/Form.svelte";
  import { z } from "zod";

  const schema = z.object({
    email: z.string().email(),
    password: z.string().min(1),
    captcha: z.string().min(1),
  });
</script>

<section class="section">
  <div class="columns is-centered">
    <div class="column is-one-third">
      <Form
        {schema}
        onSubmit={(values) => {
          console.log(values);
          router.navigate(Route.Home);
        }}
      >
        <Input
          label="Email"
          name="email"
          type="email"
          placeholder="email@example.com"
        />
        <Input
          label="Password"
          name="password"
          type="password"
          placeholder="••••••••••"
        />
        <Input
          label="reCAPTCHA Token"
          name="captcha"
          type="text"
          placeholder="Strava's reCAPTCHA token"
        />
        <div class="field">
          <div class="control">
            <button class="button is-primary is-light is-outlined px-5"
              >Login</button
            >
          </div>
        </div>
      </Form>
    </div>
  </div>
</section>
