import { createSignal, Show } from "solid-js";
import CalendarDashboard from "~/components/dashboard/CalendarDashboard";
import MusicDashboard from "~/components/dashboard/MusicDashboard";
import Button from "~/components/ui/Button";
import Icon from "~/components/ui/Icon";
import Page from "~/components/ui/Page";

interface ServiceOption {
    id: string;
    name: string;
    description: string;
    icon: "calendar-icon" | "music-icon";
}

export default function Dashboard() {
    const [selectedService, setSelectedService] = createSignal<string | null>(null);
    const [isTransitioning, setIsTransitioning] = createSignal(false);

    const serviceOptions: ServiceOption[] = [
        {
            id: "calendar",
            name: "Calendar Services",
            description: "Sync your calendars across Google, Outlook, Apple, and more",
            icon: "calendar-icon"
        },
        {
            id: "music",
            name: "Music Services",
            description: "Connect Spotify, Apple Music, YouTube Music, and other platforms",
            icon: "music-icon"
        }
    ];

    const handleServiceClick = (serviceId: string) => {
        setIsTransitioning(true);
        // Small delay to trigger fade out animation
        setTimeout(() => {
            setSelectedService(serviceId);
            setIsTransitioning(false);
        }, 200);
    };

    const handleBack = () => {
        setIsTransitioning(true);
        // Small delay to trigger fade out animation
        setTimeout(() => {
            setSelectedService(null);
            setIsTransitioning(false);
        }, 200);
    };

    return (
        <Page>
            {/* Calendar Dashboard */}
            <Show when={selectedService() === "calendar"}>
                <div class={`transition-opacity duration-300 ease-in-out ${isTransitioning() ? 'opacity-0' : 'opacity-100'
                    }`}>
                    <CalendarDashboard onBack={handleBack} />
                </div>
            </Show>

            {/* Main Dashboard */}
            <Show when={selectedService() === null}>
                <div class={`transition-opacity duration-300 ease-in-out ${isTransitioning() ? 'opacity-0' : 'opacity-100'
                    }`}>
                    <div class="container mx-auto px-6 py-16 max-w-5xl">
                        {/* Header Section */}
                        <div class="text-center mb-16">
                            <h1 class="text-3xl md:text-4xl font-bold text-gray-900 dark:text-white mb-6">
                                Connect Your Services
                            </h1>
                            <p class="text-lg text-gray-600 dark:text-gray-300 max-w-2xl mx-auto leading-relaxed">
                                Choose a service type to start synchronizing your data across platforms
                            </p>
                        </div>

                        {/* Service Options Grid */}
                        <div class="grid md:grid-cols-2 gap-8 mb-16 mx-auto w-2/3">
                            {serviceOptions.map((service) => (
                                <div
                                    class={`group relative bg-white dark:bg-gray-800 rounded-xl border-2 transition-all duration-200 cursor-pointer hover:shadow-lg ${selectedService() === service.id
                                        ? 'border-emerald-500 shadow-lg'
                                        : 'border-gray-200 dark:border-gray-700 hover:border-emerald-300 dark:hover:border-emerald-600'
                                        }`}
                                    onClick={() => handleServiceClick(service.id)}
                                >
                                    <div class="p-8">
                                        {/* Icon Container */}
                                        <div class="mb-6">
                                            <div class="w-16 h-16 rounded-lg bg-emerald-100 dark:bg-emerald-900/30 flex items-center justify-center">
                                                <Icon
                                                    name={service.icon}
                                                    class="w-8 h-8 text-emerald-600 dark:text-emerald-400"
                                                />
                                            </div>
                                        </div>

                                        {/* Service Content */}
                                        <div class="mb-6">
                                            <h2 class="text-xl font-semibold text-gray-900 dark:text-white mb-3">
                                                {service.name}
                                            </h2>
                                            <p class="text-gray-600 dark:text-gray-300 leading-relaxed">
                                                {service.description}
                                            </p>
                                        </div>

                                        {/* Action Button */}
                                        <Button
                                            variant="primary"
                                            size="md"
                                            class="w-full"
                                            onClick={() => handleServiceClick(service.id)}
                                        >
                                            Get Started
                                        </Button>
                                    </div>
                                </div>
                            ))}
                        </div>

                        {/* Info Section */}
                        <div class="text-center">
                            <div class="inline-flex items-center justify-center p-6 bg-white dark:bg-gray-800 rounded-xl border border-gray-200 dark:border-gray-700">
                                <Icon name="sync-icon" class="w-6 h-6 text-emerald-600 dark:text-emerald-400 mr-3 flex-shrink-0" />
                                <div class="text-left">
                                    <h3 class="font-medium text-gray-900 dark:text-white mb-1">
                                        Real-time Synchronization
                                    </h3>
                                    <p class="text-sm text-gray-600 dark:text-gray-300">
                                        Your data stays synchronized across all connected services automatically
                                    </p>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </Show>

            {/* Music Dashboard */}
            <Show when={selectedService() === "music"}>
                <div class={`transition-opacity duration-300 ease-in-out ${isTransitioning() ? 'opacity-0' : 'opacity-100'
                    }`}>
                    <MusicDashboard onBack={handleBack} />
                </div>
            </Show>
        </Page>
    );
}
