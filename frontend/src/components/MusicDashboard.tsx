import { createSignal } from "solid-js";
import Button from "~/components/Button";
import Icon from "~/components/Icon";
import MusicConfigModal from "./MusicConfigModal";

interface MusicConfig {
    serviceType: string;
    instanceName: string;
    syncFrequency: string;
    conflictResolution: string;
}

interface ServiceType {
    id: string;
    name: string;
    description: string;
    icon: "spotify-icon" | "apple-music-icon" | "youtube-icon" | "deezer-icon" | "tidal-icon";
}

interface MusicService extends ServiceType {
    instanceId: string;
    instanceName?: string;
    connected: boolean;
    lastSync?: string;
}

interface MusicDashboardProps {
    onBack: () => void;
}

export default function MusicDashboard(props: MusicDashboardProps) {
    const serviceTypes: ServiceType[] = [
        {
            id: "spotify",
            name: "Spotify",
            description: "Sync your playlists, liked songs, and music library",
            icon: "spotify-icon"
        },
        {
            id: "apple-music",
            name: "Apple Music",
            description: "Connect your Apple Music library and playlists",
            icon: "apple-music-icon"
        },
        {
            id: "youtube-music",
            name: "YouTube Music",
            description: "Sync your YouTube Music playlists and liked videos",
            icon: "youtube-icon"
        },
        {
            id: "deezer",
            name: "Deezer",
            description: "Sync your Deezer favorites and playlists",
            icon: "deezer-icon"
        },
        {
            id: "tidal",
            name: "Tidal",
            description: "Connect your Tidal library and high-quality music",
            icon: "tidal-icon"
        }
    ];

    const [services, setServices] = createSignal<MusicService[]>([
        {
            id: "spotify",
            instanceId: "spotify-1",
            name: "Spotify",
            instanceName: "Personal Account",
            description: "Sync your playlists, liked songs, and music library",
            icon: "spotify-icon",
            connected: true,
            lastSync: "5 minutes ago"
        }
    ]);

    const [isModalOpen, setIsModalOpen] = createSignal(false);

    const handleConnect = (instanceId: string) => {
        setServices(prev => prev.map(service =>
            service.instanceId === instanceId
                ? { ...service, connected: !service.connected, lastSync: service.connected ? undefined : "Just now" }
                : service
        ));
    };

    const addServiceInstance = (serviceTypeId: string) => {
        const serviceType = serviceTypes.find(t => t.id === serviceTypeId);
        if (!serviceType) return;

        const existingInstances = services().filter(s => s.id === serviceTypeId);
        const instanceNumber = existingInstances.length + 1;

        const newService: MusicService = {
            ...serviceType,
            instanceId: `${serviceTypeId}-${instanceNumber}`,
            instanceName: instanceNumber > 1 ? `Account ${instanceNumber}` : undefined,
            connected: false
        };

        setServices(prev => [...prev, newService]);
    };

    const handleModalSave = (config: MusicConfig) => {
        const serviceType = serviceTypes.find(t => t.id === config.serviceType);
        if (!serviceType) return;

        const existingInstances = services().filter(s => s.id === config.serviceType);
        const instanceNumber = existingInstances.length + 1;

        const newService: MusicService = {
            ...serviceType,
            instanceId: `${config.serviceType}-${instanceNumber}`,
            instanceName: config.instanceName !== serviceType.name ? config.instanceName : undefined,
            connected: false
        };

        setServices(prev => [...prev, newService]);
    };

    const removeServiceInstance = (instanceId: string) => {
        setServices(prev => prev.filter(service => service.instanceId !== instanceId));
    };

    const getServicesByType = (typeId: string) => {
        return services().filter(service => service.id === typeId);
    };

    const connectedCount = () => services().filter(s => s.connected).length;

    return (
        <div class="min-h-screen w-full bg-gray-50 dark:bg-gray-900">
            <div class="container mx-auto px-6 py-16 max-w-5xl">
                {/* Header Section with Universal Add Button */}
                <div class="flex items-center justify-between mb-8">
                    <div class="flex items-center">
                        <Button
                            variant="secondary"
                            size="sm"
                            onClick={props.onBack}
                            class="mr-4"
                        >
                            <Icon name="arrow-left" class="w-4 h-4 mr-2" />
                            Back
                        </Button>
                        <div>
                            <h1 class="text-3xl md:text-4xl font-bold text-gray-900 dark:text-white">
                                Music Integration
                            </h1>
                            <p class="text-lg text-gray-600 dark:text-gray-300 mt-2">
                                Connect and sync your music streaming services
                            </p>
                        </div>
                    </div>
                    <Button
                        variant="primary"
                        size="md"
                        onClick={() => setIsModalOpen(true)}
                    >
                        <Icon name="plus-icon" class="w-5 h-5 mr-2" />
                        Add Music Service
                    </Button>
                </div>

                {/* Status Overview */}
                <div class="grid md:grid-cols-3 gap-6 mb-12">
                    <div class="bg-white dark:bg-gray-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700">
                        <div class="flex items-center justify-between">
                            <div>
                                <p class="text-sm font-medium text-gray-600 dark:text-gray-400">Connected Services</p>
                                <p class="text-2xl font-bold text-emerald-600 dark:text-emerald-400">{connectedCount()}</p>
                            </div>
                            <Icon name="music-icon" class="w-8 h-8 text-emerald-600 dark:text-emerald-400" />
                        </div>
                    </div>

                    <div class="bg-white dark:bg-gray-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700">
                        <div class="flex items-center justify-between">
                            <div>
                                <p class="text-sm font-medium text-gray-600 dark:text-gray-400">Sync Status</p>
                                <p class="text-2xl font-bold text-green-600 dark:text-green-400">Active</p>
                            </div>
                            <Icon name="sync-icon" class="w-8 h-8 text-green-600 dark:text-green-400" />
                        </div>
                    </div>

                    <div class="bg-white dark:bg-gray-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700">
                        <div class="flex items-center justify-between">
                            <div>
                                <p class="text-sm font-medium text-gray-600 dark:text-gray-400">Last Sync</p>
                                <p class="text-2xl font-bold text-blue-600 dark:text-blue-400">5m</p>
                            </div>
                            <Icon name="clock-icon" class="w-8 h-8 text-blue-600 dark:text-blue-400" />
                        </div>
                    </div>
                </div>

                {/* Music Services */}
                <div class="mb-12">
                    <h2 class="text-2xl font-bold text-gray-900 dark:text-white mb-6">Music Services</h2>
                    <div class="space-y-8">
                        {serviceTypes.map((serviceType) => {
                            const instances = getServicesByType(serviceType.id);
                            return (
                                <div class="space-y-4">
                                    {/* Service Type Header */}
                                    <div class="flex items-center justify-between">
                                        <div class="flex items-center">
                                            <div class="w-10 h-10 rounded-lg bg-emerald-100 dark:bg-emerald-900/30 flex items-center justify-center mr-3">
                                                <Icon
                                                    name={serviceType.icon}
                                                    class="w-5 h-5 text-emerald-600 dark:text-emerald-400"
                                                />
                                            </div>
                                            <div>
                                                <h3 class="text-lg font-semibold text-gray-900 dark:text-white">
                                                    {serviceType.name}
                                                </h3>
                                                <p class="text-sm text-gray-600 dark:text-gray-300">
                                                    {serviceType.description}
                                                </p>
                                            </div>
                                        </div>
                                        <Button
                                            variant="secondary"
                                            size="sm"
                                            onClick={() => addServiceInstance(serviceType.id)}
                                        >
                                            <Icon name="plus-icon" class="w-4 h-4 mr-2" />
                                            Add Account
                                        </Button>
                                    </div>

                                    {/* Service Instances */}
                                    {instances.length > 0 ? (
                                        <div class="grid md:grid-cols-2 gap-4 ml-13">
                                            {instances.map((service) => (
                                                <div class="bg-white dark:bg-gray-800 rounded-xl border border-gray-200 dark:border-gray-700 p-6">
                                                    <div class="flex items-start justify-between mb-4">
                                                        <div>
                                                            <h4 class="text-base font-medium text-gray-900 dark:text-white">
                                                                {service.instanceName || service.name}
                                                            </h4>
                                                            {service.instanceName && (
                                                                <p class="text-sm text-gray-500 dark:text-gray-400">
                                                                    {service.name}
                                                                </p>
                                                            )}
                                                        </div>
                                                        <div class="flex items-center gap-2">
                                                            <div class={`px-2 py-1 rounded-full text-xs font-medium ${service.connected
                                                                ? 'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400'
                                                                : 'bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-300'
                                                                }`}>
                                                                {service.connected ? 'Connected' : 'Disconnected'}
                                                            </div>
                                                            <Button
                                                                variant="secondary"
                                                                size="sm"
                                                                onClick={() => removeServiceInstance(service.instanceId)}
                                                                class="p-1"
                                                            >
                                                                <Icon name="trash-icon" class="w-4 h-4" />
                                                            </Button>
                                                        </div>
                                                    </div>

                                                    {service.lastSync && (
                                                        <p class="text-xs text-gray-500 dark:text-gray-400 mb-4">
                                                            Last synced: {service.lastSync}
                                                        </p>
                                                    )}

                                                    <div class="flex gap-3">
                                                        <Button
                                                            variant={service.connected ? "secondary" : "primary"}
                                                            size="sm"
                                                            onClick={() => handleConnect(service.instanceId)}
                                                            class="flex-1"
                                                        >
                                                            {service.connected ? 'Disconnect' : 'Connect'}
                                                        </Button>
                                                        {service.connected && (
                                                            <Button
                                                                variant="secondary"
                                                                size="sm"
                                                                onClick={() => console.log(`Configure ${service.instanceId}`)}
                                                            >
                                                                Configure
                                                            </Button>
                                                        )}
                                                    </div>
                                                </div>
                                            ))}
                                        </div>
                                    ) : (
                                        <div class="ml-13 p-4 border-2 border-dashed border-gray-300 dark:border-gray-600 rounded-lg text-center">
                                            <p class="text-gray-500 dark:text-gray-400 text-sm">
                                                No accounts connected. Click "Add Account" to get started.
                                            </p>
                                        </div>
                                    )}
                                </div>
                            );
                        })}
                    </div>
                </div>

                {/* Sync Settings */}
                <div class="bg-white dark:bg-gray-800 rounded-xl border border-gray-200 dark:border-gray-700 p-6">
                    <h3 class="text-xl font-semibold text-gray-900 dark:text-white mb-4">Sync Settings</h3>
                    <div class="grid md:grid-cols-2 gap-6">
                        <div>
                            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                                Sync Frequency
                            </label>
                            <select
                                class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-white"
                                aria-label="Sync Frequency"
                            >
                                <option>Every 5 minutes</option>
                                <option>Every 15 minutes</option>
                                <option>Every hour</option>
                                <option>Manual only</option>
                            </select>
                        </div>
                        <div>
                            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                                Conflict Resolution
                            </label>
                            <select
                                class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-white"
                                aria-label="Conflict Resolution"
                            >
                                <option>Keep both</option>
                                <option>Merge when possible</option>
                            </select>
                        </div>
                    </div>
                    <div class="mt-6">
                        <Button variant="primary" size="md">
                            Save Settings
                        </Button>
                    </div>
                </div>
            </div>

            {/* Music Configuration Modal */}
            <MusicConfigModal
                isOpen={isModalOpen}
                onClose={() => setIsModalOpen(false)}
                onSave={handleModalSave}
            />
        </div>
    );
} 