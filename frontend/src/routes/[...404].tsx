import { A } from "@solidjs/router";

export default function NotFound() {
  return (
    <main class="text-center mx-auto text-slate-900 dark:text-slate-100 p-4">
      <h1 class="max-6-xs text-6xl text-emerald-600 dark:text-emerald-400 font-thin uppercase my-16">Not Found</h1>
      <p class="mt-8">
        Visit{" "}
        <a href="https://solidjs.com" target="_blank" rel="noopener" class="text-emerald-600 dark:text-emerald-400 hover:underline">
          solidjs.com
        </a>{" "}
        to learn how to build Solid apps.
      </p>
      <p class="my-4">
        <A href="/" class="text-emerald-600 dark:text-emerald-400 hover:underline">
          Home
        </A>
        {" - "}
        <A href="/about" class="text-emerald-600 dark:text-emerald-400 hover:underline">
          About Page
        </A>
      </p>
    </main>
  );
}
