import { createSignal } from "solid-js";
import { BaseDashboard } from "~/components/dashboard";
import type { ExtendedServiceInstanceSyncSettings, ServiceInstance, ServiceSpecificConfig, ServiceType } from "~/types";
import {
    CheckboxField,
    ConfigSection,
    SelectField,
    ServiceInstanceConfigModal,
} from ".";
import MusicConfigModal from "../modals/MusicConfigModal";

interface MusicConfig {
    serviceType: string;
    instanceName: string;
    syncFrequency: string;
    conflictResolution: string;
}

interface MusicService extends ServiceInstance {
    id: string;
    description: string;
    icon: "spotify-icon" | "apple-music-icon" | "youtube-icon" | "deezer-icon" | "tidal-icon";
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
            lastSync: "5 minutes ago",
            syncSettings: {
                frequency: "Every 30 minutes",
                conflictResolution: "Merge when possible"
            }
        }
    ]);

    const [isAddModalOpen, setIsAddModalOpen] = createSignal(false);
    const [isConfigModalOpen, setIsConfigModalOpen] = createSignal(false);
    const [selectedServiceForConfig, setSelectedServiceForConfig] = createSignal<MusicService | null>(null);

    // Music-specific configuration signals
    const [syncPlaylists, setSyncPlaylists] = createSignal<boolean>(true);
    const [syncLikedSongs, setSyncLikedSongs] = createSignal<boolean>(true);
    const [syncLibrary, setSyncLibrary] = createSignal<boolean>(true);
    const [qualityPreference, setQualityPreference] = createSignal<string>("High");
    const [bidirectionalSync, setBidirectionalSync] = createSignal<boolean>(true);
    const [autoResolveConflicts, setAutoResolveConflicts] = createSignal<boolean>(false);
    const [privateMode, setPrivateMode] = createSignal<boolean>(false);

    const handleConnect = (instanceId: string) => {
        setServices(prev => prev.map(service =>
            service.instanceId === instanceId
                ? { ...service, connected: !service.connected, lastSync: service.connected ? undefined : "Just now" }
                : service
        ));
    };

    const handleConfigure = (instanceId: string) => {
        const service = services().find(s => s.instanceId === instanceId);
        if (service) {
            setSelectedServiceForConfig(service);
            setIsConfigModalOpen(true);
        }
    };

    const handleConfigSave = (instanceId: string, instanceName: string, syncSettings: ExtendedServiceInstanceSyncSettings) => {
        setServices(prev => prev.map(service =>
            service.instanceId === instanceId
                ? {
                    ...service,
                    instanceName: instanceName !== service.name ? instanceName : undefined,
                    syncSettings
                }
                : service
        ));
        setSelectedServiceForConfig(null);
    };

    const handleSync = (instanceId: string) => {
        console.log("Syncing instance", instanceId);
    };

    const addServiceInstance = (serviceTypeId: string) => {
        const serviceType = serviceTypes.find(t => t.id === serviceTypeId);
        if (!serviceType) return;

        const existingInstances = services().filter(s => s.id === serviceTypeId);
        const instanceNumber = existingInstances.length + 1;

        const newService: MusicService = {
            id: serviceType.id,
            name: serviceType.name,
            description: serviceType.description,
            icon: serviceType.icon as MusicService["icon"],
            instanceId: `${serviceTypeId}-${instanceNumber}`,
            instanceName: instanceNumber > 1 ? `Account ${instanceNumber}` : undefined,
            connected: false,
            syncSettings: {
                frequency: "Every 15 minutes",
                conflictResolution: "Keep both"
            }
        };

        setServices(prev => [...prev, newService]);
    };

    const handleModalSave = (config: MusicConfig) => {
        const serviceType = serviceTypes.find(t => t.id === config.serviceType);
        if (!serviceType) return;

        const existingInstances = services().filter(s => s.id === config.serviceType);
        const instanceNumber = existingInstances.length + 1;

        const newService: MusicService = {
            id: serviceType.id,
            name: serviceType.name,
            description: serviceType.description,
            icon: serviceType.icon as MusicService["icon"],
            instanceId: `${config.serviceType}-${instanceNumber}`,
            instanceName: config.instanceName !== serviceType.name ? config.instanceName : undefined,
            connected: false,
            syncSettings: {
                frequency: config.syncFrequency,
                conflictResolution: config.conflictResolution
            }
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

    // Music-specific configuration handlers
    const getServiceSpecificConfig = (): ServiceSpecificConfig => {
        return {
            syncPlaylists: syncPlaylists(),
            syncLikedSongs: syncLikedSongs(),
            syncLibrary: syncLibrary(),
            qualityPreference: qualityPreference(),
            bidirectionalSync: bidirectionalSync(),
            autoResolveConflicts: autoResolveConflicts(),
            privateMode: privateMode()
        };
    };

    const initializeServiceSpecificConfig = (config: ServiceSpecificConfig | undefined) => {
        setSyncPlaylists(config?.syncPlaylists ?? true);
        setSyncLikedSongs(config?.syncLikedSongs ?? true);
        setSyncLibrary(config?.syncLibrary ?? true);
        setQualityPreference(config?.qualityPreference || "High");
        setBidirectionalSync(config?.bidirectionalSync ?? true);
        setAutoResolveConflicts(config?.autoResolveConflicts ?? false);
        setPrivateMode(config?.privateMode ?? false);
    };

    const resetServiceSpecificConfig = () => {
        setSyncPlaylists(true);
        setSyncLikedSongs(true);
        setSyncLibrary(true);
        setQualityPreference("High");
        setBidirectionalSync(true);
        setAutoResolveConflicts(false);
        setPrivateMode(false);
    };

    // Music-specific configuration section
    const musicConfigSection = (
        <>
            <ConfigSection title="Music Sync Options">
                <div class="space-y-3">
                    <CheckboxField
                        label="Sync Playlists"
                        checked={syncPlaylists()}
                        onChange={setSyncPlaylists}
                    />
                    <CheckboxField
                        label="Sync Liked Songs"
                        checked={syncLikedSongs()}
                        onChange={setSyncLikedSongs}
                    />
                    <CheckboxField
                        label="Sync Music Library"
                        checked={syncLibrary()}
                        onChange={setSyncLibrary}
                    />
                </div>

                <SelectField
                    label="Quality Preference"
                    value={qualityPreference()}
                    onChange={setQualityPreference}
                    ariaLabel="Quality Preference"
                    options={[
                        { value: "Low", label: "Low" },
                        { value: "Normal", label: "Normal" },
                        { value: "High", label: "High" },
                        { value: "Lossless", label: "Lossless" }
                    ]}
                />
            </ConfigSection>

            <ConfigSection title="Advanced Options">
                <div class="space-y-3">
                    <CheckboxField
                        label="Bidirectional Sync"
                        description="Sync changes in both directions"
                        checked={bidirectionalSync()}
                        onChange={setBidirectionalSync}
                    />
                    <CheckboxField
                        label="Auto-resolve Conflicts"
                        description="Automatically apply conflict resolution rules"
                        checked={autoResolveConflicts()}
                        onChange={setAutoResolveConflicts}
                    />
                    <CheckboxField
                        label="Private Mode"
                        description="Hide personal data during sync operations"
                        checked={privateMode()}
                        onChange={setPrivateMode}
                    />
                </div>
            </ConfigSection>
        </>
    );

    return (
        <BaseDashboard
            title="Music Integration"
            description="Connect and sync your music streaming services"
            addButtonText="Add Music Service"
            mainIcon="music-icon"
            serviceTypes={serviceTypes}
            services={services}
            connectedCount={connectedCount}
            lastSyncValue="5m"
            onBack={props.onBack}
            onAdd={() => setIsAddModalOpen(true)}
            onAddAccount={addServiceInstance}
            onConnect={handleConnect}
            onSync={handleSync}
            onRemove={removeServiceInstance}
            onConfigure={handleConfigure}
            getServicesByType={getServicesByType}
            modal={
                <>
                    <MusicConfigModal
                        isOpen={isAddModalOpen}
                        onClose={() => setIsAddModalOpen(false)}
                        onSave={handleModalSave}
                    />
                    <ServiceInstanceConfigModal
                        isOpen={isConfigModalOpen}
                        onClose={() => {
                            setIsConfigModalOpen(false);
                            setSelectedServiceForConfig(null);
                        }}
                        onSave={handleConfigSave}
                        serviceInstance={selectedServiceForConfig}
                        customConfigSection={musicConfigSection}
                        getServiceSpecificConfig={getServiceSpecificConfig}
                        initializeServiceSpecificConfig={initializeServiceSpecificConfig}
                        resetServiceSpecificConfig={resetServiceSpecificConfig}
                    />
                </>
            }
        />
    );
} 