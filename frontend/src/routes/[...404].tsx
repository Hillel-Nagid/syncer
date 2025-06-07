import { A } from "@solidjs/router";

export default function NotFound() {
  return (
    <main class="min-h-screen flex flex-col items-center justify-center text-center mx-auto text-slate-900 dark:text-slate-100 p-4">
      <div class="max-w-md w-full">
        {/* Large 404 number */}
        <h1 class="text-8xl font-bold text-emerald-600 dark:text-emerald-400 mb-4">
          404
        </h1>

        {/* Main heading */}
        <h2 class="text-3xl font-semibold mb-6">
          Page Not Found
        </h2>

        {/* Description */}
        <p class="text-lg text-slate-600 dark:text-slate-400 mb-8">
          The sync endpoint you're looking for doesn't exist or has been moved.
          Let's get you back to managing your synchronizations.
        </p>

        {/* Action buttons */}
        <div class="space-y-4">
          <A
            href="/"
            class="inline-block bg-emerald-600 hover:bg-emerald-700 text-white font-medium py-3 px-6 rounded-lg transition-colors duration-200"
          >
            Back to Home Page
          </A>
        </div>
      </div>
    </main>
  );
}
