import { A } from "@solidjs/router";
import Button from "~/components/ui/Button";
import Card from "~/components/ui/Card";
import Container from "~/components/ui/Container";
import FeatureCard from "~/components/ui/FeatureCard";
import Hero from "~/components/ui/Hero";
import Icon from "~/components/ui/Icon";
import Page from "~/components/ui/Page";
import type { IconName } from "~/types";

// Hero Section Component
function HeroSection() {
  return (
    <section class="relative bg-gradient-to-br from-emerald-50 to-emerald-100 dark:from-slate-900 dark:to-slate-800 py-20 px-4" >
      <Container maxWidth="xl">
        <Hero
          title="Syncer"
          subtitle="Seamlessly synchronize your digital life. Connect calendars, music streaming services, and more with our powerful synchronization engine."
        />
        <div class="flex flex-col sm:flex-row gap-4 justify-center">
          <A href="/dashboard">
            <Button variant="primary" size="xl">
              Get Started
            </Button>
          </A>
        </div>
      </Container>
    </section>
  );
}

// Features Section Component
function FeaturesSection() {
  return (
    <section class="py-20">
      <Container>
        <h2 class="text-3xl md:text-4xl font-bold text-center mb-12">
          Sync Across Platforms
        </h2>
        <div class="grid md:grid-cols-2 lg:grid-cols-3 gap-8">
          <FeatureCard
            icon="calendar-icon"
            title="Calendar Sync"
            description="Keep all your calendars in perfect sync across Google Calendar, Outlook, Apple Calendar, and more."
          />
          <FeatureCard
            icon="music-icon"
            title="Music Streaming"
            description="Synchronize playlists, favorites, and listening history across Spotify, Apple Music, YouTube Music, and more."
          />
          <FeatureCard
            icon="realtime-sync-icon"
            title="Real-time Sync"
            description="Changes are synchronized instantly across all your connected services with our real-time sync engine."
          />
        </div>
      </Container>
    </section>
  );
}

// Supported Services Section Component
function SupportedServicesSection() {
  const services: { name: string; icon: IconName; description: string }[] = [
    {
      name: "Google Calendar",
      icon: "google-calendar-icon" as const,
      description: "Sync events and schedules"
    },
    {
      name: "Spotify",
      icon: "spotify-icon" as const,
      description: "Sync playlists and music"
    },
    {
      name: "Deezer",
      icon: "deezer-icon" as const,
      description: "Sync music library"
    },
    {
      name: "Apple Music",
      icon: "apple-music-icon" as const,
      description: "Sync music and playlists"
    },
    {
      name: "YouTube Music",
      icon: "youtube-icon" as const,
      description: "Sync music and subscriptions"
    },
    {
      name: "Outlook Calendar",
      icon: "outlook-icon" as const,
      description: "Sync events and meetings"
    },
    {
      name: "Tidal",
      icon: "tidal-icon" as const,
      description: "Sync high-quality music"
    }
  ];

  return (
    <section class="py-20 px-4 bg-slate-50 dark:bg-slate-800/50" >
      <Container>
        <div class="text-center">
          <h2 class="text-3xl md:text-4xl font-bold mb-12">
            Supported Services
          </h2>
          <div class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
            {services.map((service) => (
              <Card padding="md" class="group hover:shadow-xl hover:border-emerald-300 dark:hover:border-emerald-600 transition-all duration-300 hover:scale-105">
                <div class="flex flex-col items-center text-center space-y-4">
                  <div class="w-16 h-16 flex items-center justify-center bg-gradient-to-br from-emerald-50 to-emerald-100 dark:from-emerald-900/20 dark:to-emerald-800/20 rounded-2xl group-hover:from-emerald-100 group-hover:to-emerald-200 dark:group-hover:from-emerald-800/30 dark:group-hover:to-emerald-700/30 transition-all duration-300">
                    <Icon name={service.icon} class="w-8 h-8" alt={`${service.name} Logo`} />
                  </div>
                  <div>
                    <h3 class="font-semibold text-slate-900 dark:text-slate-100 mb-1">
                      {service.name}
                    </h3>
                    <p class="text-sm text-slate-600 dark:text-slate-400">
                      {service.description}
                    </p>
                  </div>
                </div>
              </Card>
            ))}
          </div>
        </div>
      </Container>
    </section>
  );
}

// Call to Action Section Component
function CallToActionSection() {
  return (
    <section class="py-20">
      <Container maxWidth="lg">
        <div class="text-center">
          <h2 class="text-3xl md:text-4xl font-bold mb-6">
            Ready to sync your digital life?
          </h2>
          <p class="text-xl text-slate-600 dark:text-slate-400 mb-8">
            Join thousands of users who have simplified their digital workflow with Syncer.
          </p>
          <A href="/dashboard">
            <Button variant="primary" size="lg" class="mx-auto justify-self-center">
              Start Syncing Now
            </Button>
          </A>
        </div>
      </Container>
    </section>
  );
}



// Main Home Component
export default function Home() {
  return (
    <Page>
      <main class="text-slate-900 dark:text-slate-100">
        <HeroSection />
        <FeaturesSection />
        <SupportedServicesSection />
        <CallToActionSection />
      </main>
    </Page>
  );
}
