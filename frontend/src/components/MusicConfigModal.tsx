import { createSignal, For } from "solid-js";
import Button from "./Button";
import Icon from "./Icon";
import Modal from "./Modal";

interface ServiceType {
    id: string;
    name: string;
    description: string;
    icon: "spotify-icon" | "apple-music-icon" | "youtube-icon" | "deezer-icon" | "tidal-icon";
}

interface MusicConfigModalProps {
    isOpen: () => boolean;
    onClose: () => void;
    onSave: (config: MusicConfig) => void;
}

interface MusicConfig {
    serviceType: string;
    instanceName: string;
    syncFrequency: string;
    conflictResolution: string;
}

export default function MusicConfigModal(props: MusicConfigModalProps) {
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

    const [selectedServiceType, setSelectedServiceType] = createSignal<string>("");
    const [instanceName, setInstanceName] = createSignal<string>("");
    const [syncFrequency, setSyncFrequency] = createSignal<string>("Every 15 minutes");
    const [conflictResolution, setConflictResolution] = createSignal<string>("Keep both");

    const resetForm = () => {
        setSelectedServiceType("");
        setInstanceName("");
        setSyncFrequency("Every 15 minutes");
        setConflictResolution("Keep both");
    };

    const handleClose = () => {
        resetForm();
        props.onClose();
    };

    const handleSave = () => {
        if (!selectedServiceType()) return;

        const config: MusicConfig = {
            serviceType: selectedServiceType(),
            instanceName: instanceName() || getSelectedServiceName(),
            syncFrequency: syncFrequency(),
            conflictResolution: conflictResolution()
        };

        props.onSave(config);
        resetForm();
        props.onClose();
    };

    const getSelectedServiceName = () => {
        const service = serviceTypes.find(s => s.id === selectedServiceType());
        return service?.name || "";
    };

    const isFormValid = () => selectedServiceType() !== "";

    return (
        <Modal
            isOpen={props.isOpen}
            onClose={handleClose}
            title="Add Music Service"
            size="lg"
        >
            <div class="space-y-6">
                {/* Service Type Selection */}
                <div>
                    <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-3">
                        Select Music Streaming Service
                    </label>
                    <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
                        <For each={serviceTypes}>
                            {(service) => (
                                <div
                                    class={`relative cursor-pointer rounded-lg border-2 p-4 transition-all ${selectedServiceType() === service.id
                                        ? 'border-emerald-500 bg-emerald-50 dark:bg-emerald-900/20'
                                        : 'border-gray-200 dark:border-gray-600 hover:border-emerald-300 dark:hover:border-emerald-600'
                                        }`}
                                    onClick={() => setSelectedServiceType(service.id)}
                                >
                                    <div class="flex items-center">
                                        <div class="flex-shrink-0">
                                            <div class="w-10 h-10 rounded-lg bg-emerald-100 dark:bg-emerald-900/30 flex items-center justify-center">
                                                <Icon
                                                    name={service.icon}
                                                    class="w-5 h-5 text-emerald-600 dark:text-emerald-400"
                                                />
                                            </div>
                                        </div>
                                        <div class="ml-3 min-w-0 flex-1">
                                            <p class="text-sm font-medium text-gray-900 dark:text-white">
                                                {service.name}
                                            </p>
                                            <p class="text-xs text-gray-500 dark:text-gray-400 truncate">
                                                {service.description}
                                            </p>
                                        </div>
                                        {selectedServiceType() === service.id && (
                                            <div class="flex-shrink-0">
                                                <Icon
                                                    name="v-icon"
                                                    class="w-5 h-5 text-emerald-600 dark:text-emerald-400"
                                                />
                                            </div>
                                        )}
                                    </div>
                                </div>
                            )}
                        </For>
                    </div>
                </div>

                {/* Instance Name */}
                <div>
                    <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                        Account Name (Optional)
                    </label>
                    <input
                        type="text"
                        value={instanceName()}
                        onInput={(e) => setInstanceName(e.currentTarget.value)}
                        placeholder={`${getSelectedServiceName()} Account`}
                        class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-white placeholder-gray-500 dark:placeholder-gray-400 focus:ring-2 focus:ring-emerald-500 focus:border-emerald-500"
                    />
                    <p class="text-xs text-gray-500 dark:text-gray-400 mt-1">
                        Give this music service a custom name to help identify it
                    </p>
                </div>

                {/* Configuration Options */}
                <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
                    <div>
                        <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                            Sync Frequency
                        </label>
                        <select
                            value={syncFrequency()}
                            onChange={(e) => setSyncFrequency(e.currentTarget.value)}
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
                            value={conflictResolution()}
                            onChange={(e) => setConflictResolution(e.currentTarget.value)}
                            class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-white"
                            aria-label="Conflict Resolution"
                        >
                            <option>Keep both</option>
                            <option>Merge when possible</option>
                        </select>
                    </div>
                </div>

                {/* Music-specific options */}
                <div>
                    <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-3">
                        Sync Options
                    </label>
                    <div class="space-y-3">
                        <label class="flex items-center">
                            <input
                                type="checkbox"
                                class="w-4 h-4 text-emerald-600 border-gray-300 rounded focus:ring-emerald-500 dark:border-gray-600 dark:bg-gray-700"
                                checked
                            />
                            <span class="ml-2 text-sm text-gray-700 dark:text-gray-300">Sync playlists</span>
                        </label>
                        <label class="flex items-center">
                            <input
                                type="checkbox"
                                class="w-4 h-4 text-emerald-600 border-gray-300 rounded focus:ring-emerald-500 dark:border-gray-600 dark:bg-gray-700"
                                checked
                            />
                            <span class="ml-2 text-sm text-gray-700 dark:text-gray-300">Sync liked/favorite songs</span>
                        </label>
                        <label class="flex items-center">
                            <input
                                type="checkbox"
                                class="w-4 h-4 text-emerald-600 border-gray-300 rounded focus:ring-emerald-500 dark:border-gray-600 dark:bg-gray-700"
                            />
                            <span class="ml-2 text-sm text-gray-700 dark:text-gray-300">Sync listening history</span>
                        </label>
                        <label class="flex items-center">
                            <input
                                type="checkbox"
                                class="w-4 h-4 text-emerald-600 border-gray-300 rounded focus:ring-emerald-500 dark:border-gray-600 dark:bg-gray-700"
                            />
                            <span class="ml-2 text-sm text-gray-700 dark:text-gray-300">Sync followed artists</span>
                        </label>
                    </div>
                </div>

                {/* Action Buttons */}
                <div class="flex justify-end space-x-3 pt-4 border-t border-gray-200 dark:border-gray-700">
                    <Button variant="secondary" onClick={handleClose}>
                        Cancel
                    </Button>
                    <Button
                        variant="primary"
                        onClick={handleSave}
                        class={!isFormValid() ? "opacity-50 cursor-not-allowed" : ""}
                    >
                        <Icon name="plus-icon" class="w-4 h-4 mr-2" />
                        Add Music Service
                    </Button>
                </div>
            </div>
        </Modal>
    );
} 