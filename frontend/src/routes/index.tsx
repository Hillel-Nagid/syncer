import { A } from "@solidjs/router";
import { Show } from "solid-js";
import Button from "~/components/Button";
import Icon from "~/components/Icon";
import { useUser } from "~/contexts/UserContext";

export default function Home() {
  const { isLoggedIn, login, logout, user } = useUser();

  const handleDemoLogin = () => {
    login({
      id: "demo-user",
      username: "John Doe",
      email: "john@example.com",
      profilePicture: "https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=150&h=150&fit=crop&crop=face"
    });
  };
  const services = [
    "Google Calendar",
    "Spotify",
    "Deezer"
  ]
  return (
    <main class="text-slate-900 dark:text-slate-100">
      <section class="relative bg-gradient-to-br from-emerald-50 to-emerald-100 dark:from-slate-900 dark:to-slate-800 py-20 px-4">
        <div class="max-w-6xl mx-auto text-center">
          <h1 class="text-5xl md:text-6xl font-bold text-emerald-600 dark:text-emerald-400 mb-6">
            Syncer
          </h1>
          <p class="text-xl md:text-2xl text-slate-700 dark:text-slate-300 mb-8 max-w-3xl mx-auto">
            Seamlessly synchronize your digital life. Connect calendars, music streaming services, and more with our powerful synchronization engine.
          </p>
          <div class="flex flex-col sm:flex-row gap-4 justify-center">
            <Button variant="primary" size="xl">
              <A href="/dashboard">
                Get Started
              </A>
            </Button>
          </div>
        </div>
      </section>

      <section class="py-20 px-4">
        <div class="max-w-6xl mx-auto">
          <h2 class="text-3xl md:text-4xl font-bold text-center mb-12">
            Sync Across Platforms
          </h2>
          <div class="grid md:grid-cols-2 lg:grid-cols-3 gap-8">
            <div class="bg-white dark:bg-slate-800 p-6 rounded-xl shadow-lg border border-slate-200 dark:border-slate-700">
              <div class="w-12 h-12 bg-emerald-100 dark:bg-emerald-900/30 rounded-lg flex items-center justify-center mb-4">
                <Icon name="calendar-icon" class="w-6 h-6 text-emerald-600 dark:text-emerald-400" alt="Calendar Icon" />
              </div>
              <h3 class="text-xl font-semibold mb-3">Calendar Sync</h3>
              <p class="text-slate-600 dark:text-slate-400">
                Keep all your calendars in perfect sync across Google Calendar, Outlook, Apple Calendar, and more.
              </p>
            </div>

            <div class="bg-white dark:bg-slate-800 p-6 rounded-xl shadow-lg border border-slate-200 dark:border-slate-700">
              <div class="w-12 h-12 bg-emerald-100 dark:bg-emerald-900/30 rounded-lg flex items-center justify-center mb-4">
                <Icon name="music-icon" class="w-6 h-6 text-emerald-600 dark:text-emerald-400" alt="Music Icon" />
              </div>
              <h3 class="text-xl font-semibold mb-3">Music Streaming</h3>
              <p class="text-slate-600 dark:text-slate-400">
                Synchronize playlists, favorites, and listening history across Spotify, Apple Music, YouTube Music, and more.
              </p>
            </div>

            <div class="bg-white dark:bg-slate-800 p-6 rounded-xl shadow-lg border border-slate-200 dark:border-slate-700">
              <div class="w-12 h-12 bg-emerald-100 dark:bg-emerald-900/30 rounded-lg flex items-center justify-center mb-4">
                <Icon name="realtime-sync-icon" class="w-6 h-6 text-emerald-600 dark:text-emerald-400" alt="Real-time Sync Icon" />
              </div>
              <h3 class="text-xl font-semibold mb-3">Real-time Sync</h3>
              <p class="text-slate-600 dark:text-slate-400">
                Changes are synchronized instantly across all your connected services with our real-time sync engine.
              </p>
            </div>
          </div>
        </div>
      </section>

      <section class="py-20 px-4 bg-slate-50 dark:bg-slate-800/50">
        <div class="max-w-6xl mx-auto text-center">
          <h2 class="text-3xl md:text-4xl font-bold mb-12">
            Supported Services
          </h2>
          <div class="grid grid-cols-3 md:grid-cols-3 lg:grid-cols-3 gap-8 items-center">
            {services.map((service) => (
              <div class="bg-white dark:bg-slate-800 p-4 rounded-lg shadow-md border border-slate-200 dark:border-slate-700">
                <span class="text-sm font-medium text-slate-700 dark:text-slate-300">{service}</span>
              </div>
            ))}
          </div>
        </div>
      </section>

      <section class="py-20 px-4 bg-slate-50 dark:bg-slate-800/50">
        <div class="max-w-4xl mx-auto text-center">
          <h2 class="text-3xl md:text-4xl font-bold mb-6">
            Demo Authentication
          </h2>
          <p class="text-xl text-slate-600 dark:text-slate-400 mb-8">
            Test the authentication functionality by logging in and out.
          </p>
          <Show
            when={isLoggedIn()}
            fallback={
              <Button variant="primary" size="lg" onClick={handleDemoLogin} class="justify-self-center">
                Demo Login
              </Button>
            }
          >
            <div class="flex flex-col items-center space-y-4">
              <p class="text-lg text-slate-700 dark:text-slate-300">
                Welcome back, {user()?.username}!
              </p>
              <Button variant="secondary" size="lg" onClick={logout} class="justify-self-center">
                Demo Logout
              </Button>
            </div>
          </Show>
        </div>
      </section>

      <section class="py-20 px-4">
        <div class="max-w-4xl mx-auto text-center">
          <h2 class="text-3xl md:text-4xl font-bold mb-6">
            Ready to sync your digital life?
          </h2>
          <p class="text-xl text-slate-600 dark:text-slate-400 mb-8">
            Join thousands of users who have simplified their digital workflow with Syncer.
          </p>
          <Button variant="primary" size="lg" class="justify-self-center-safe">
            Start Syncing Now
          </Button>
        </div>
      </section>
    </main>
  );
}
