import type { Subscriber, Readable } from "svelte/store";
import type { SvelteComponent } from "svelte";
import { readable } from "svelte/store";
import page from "page";

import Home from "@pages/Home.svelte";
import Login from "@pages/Login.svelte";

export enum Route {
  Home = "/",
  Login = "/login",
}

// setRoute updates the global route store.
// Should be initialized once.
let setRoute: Subscriber<typeof SvelteComponent>;

function updateComponent(component: typeof SvelteComponent): PageJS.Callback {
  return () => {
    setRoute(component);
  };
}

const route = readable<typeof SvelteComponent>(
  Home,
  (set: Subscriber<typeof SvelteComponent>) => {
    setRoute = set;

    page(Route.Home, updateComponent(Home));
    page(Route.Login, updateComponent(Login));

    page.start();
    return () => {
      page.stop();
    };
  }
);

const router = {
  subscribe: route.subscribe,

  navigate(path: string) {
    page.show(path);
  },
};
export default router;
