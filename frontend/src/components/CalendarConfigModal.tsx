import { createSignal, For } from "solid-js";
import Button from "./Button";
import Icon from "./Icon";
import Modal from "./Modal";

interface ServiceType {
    id: string;
    name: string;
    description: string;
    icon: "google-calendar-icon" | "outlook-icon";
}

interface CalendarConfigModalProps {
    isOpen: () => boolean;
    onClose: () => void;
    onSave: (config: CalendarConfig) => void;
}

interface CalendarConfig {
    serviceType: string;
    instanceName: string;
    syncFrequency: string;
    conflictResolution: string;
}

export default function CalendarConfigModal(props: CalendarConfigModalProps) {
    const serviceTypes: ServiceType[] = [
        {
            id: "google",
            name: "Google Calendar",
            description: "Sync with your Google Calendar events and meetings",
            icon: "google-calendar-icon"
        },
        {
            id: "outlook",
            name: "Microsoft Outlook",
            description: "Connect your Outlook calendar and Exchange events",
            icon: "outlook-icon"
        },
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

        const config: CalendarConfig = {
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
            title="Add Calendar Service"
            size="lg"
        >
            <div class="space-y-6">
                {/* Service Type Selection */}
                <div>
                    <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-3">
                        Select Calendar Service
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
                        Give this calendar a custom name to help identify it
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
                        Add Calendar
                    </Button>
                </div>
            </div>
        </Modal>
    );
} 