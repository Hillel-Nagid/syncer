import { A } from "@solidjs/router";
import Page from "~/components/ui/Page";

export default function NotFound() {
  return (
    <Page class="flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8 text-center">
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
    </Page>
  );
}
