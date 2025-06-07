import { Component } from "solid-js";
import Card from "~/components/ui/Card";
import Container from "~/components/ui/Container";
import FeatureCard from "~/components/ui/FeatureCard";
import Hero from "~/components/ui/Hero";
import Page from "~/components/ui/Page";
import StepCard from "~/components/ui/StepCard";

const About: Component = () => {
    const features = [
        {
            icon: "calendar-icon" as const,
            title: "Calendar Synchronization",
            description: "Keep all your calendars in perfect sync. Google Calendar, Outlook, Apple Calendar - all working together seamlessly."
        },
        {
            icon: "music-icon" as const,
            title: "Music Streaming Sync",
            description: "Synchronize your playlists, favorites, and listening history across Spotify, Apple Music, YouTube Music, and more."
        },
        {
            icon: "lightning-icon" as const,
            title: "Real-time Updates",
            description: "Changes propagate instantly across all connected services. Add an event in one calendar, see it everywhere."
        },
        {
            icon: "lock-icon" as const,
            title: "Privacy First",
            description: "Your data is encrypted and secure. We only access what's necessary for synchronization and never store sensitive information."
        }
    ];

    const steps = [
        {
            title: "Connect Your Services",
            description: "Securely link your calendars, music streaming accounts, and other services through our encrypted connection system."
        },
        {
            title: "Configure Sync Rules",
            description: "Set up intelligent rules for how your data should be synchronized. Choose what to sync, when, and in which direction."
        },
        {
            title: "Enjoy Seamless Sync",
            description: "Sit back and let Syncer handle the rest. Your data stays in perfect harmony across all your connected services."
        }
    ];

    return (
        <Page class="py-12 px-4 sm:px-6 lg:px-8">
            <Container>
                {/* Header Section */}
                <Hero
                    title="About Syncer"
                    subtitle="Your unified synchronization engine for seamlessly connecting calendars, music streaming services, and more."
                />

                {/* Mission Section */}
                <section class="py-12">
                    <Card>
                        <h2 class="text-3xl font-semibold text-emerald-800 dark:text-emerald-200 mb-6">Our Mission</h2>
                        <p class="text-lg text-gray-700 dark:text-gray-300 leading-relaxed">
                            In today's digital world, we use multiple services for different aspects of our lives.
                            Syncer bridges the gaps between these services, creating a unified ecosystem where your
                            data flows seamlessly between platforms. Whether it's syncing your work calendar with
                            your personal schedule or keeping your music playlists updated across streaming services,
                            Syncer makes it effortless.
                        </p>
                    </Card>
                </section>

                {/* Features Grid */}
                <section class="py-12">
                    <div class="grid md:grid-cols-2 gap-8">
                        {features.map((feature) => (
                            <FeatureCard
                                icon={feature.icon}
                                title={feature.title}
                                description={feature.description}
                            />
                        ))}
                    </div>
                </section>

                {/* How It Works Section */}
                <section class="py-12">
                    <Card>
                        <h2 class="text-3xl font-semibold text-emerald-800 dark:text-emerald-200 mb-8 text-center">How It Works</h2>
                        <div class="grid md:grid-cols-3 gap-8">
                            {steps.map((step, index) => (
                                <StepCard
                                    step={index + 1}
                                    title={step.title}
                                    description={step.description}
                                />
                            ))}
                        </div>
                    </Card>
                </section>
            </Container>
        </Page>
    );
};

export default About;
