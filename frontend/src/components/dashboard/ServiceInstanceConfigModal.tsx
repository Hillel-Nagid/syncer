import { createSignal, JSX, Show } from "solid-js";
import Modal from "../modals/Modal";
import Button from "../ui/Button";
import Icon from "../ui/Icon";
import type { ServiceInstance, ServiceInstanceSyncSettings } from "./ServiceInstanceCard";

// Generic service-specific configuration type
export interface ServiceSpecificConfig {
    [key: string]: any;
}

export interface ExtendedServiceInstanceSyncSettings extends ServiceInstanceSyncSettings {
    serviceSpecific?: ServiceSpecificConfig;
}

// Reusable form primitives
export interface CheckboxFieldProps {
    label: string;
    checked: boolean;
    onChange: (checked: boolean) => void;
    description?: string;
}

export function CheckboxField(props: CheckboxFieldProps) {
    return (
        <label class="flex items-start">
            <input
                type="checkbox"
                checked={props.checked}
                onChange={(e) => props.onChange(e.currentTarget.checked)}
                class="rounded border-gray-300 text-emerald-600 focus:ring-emerald-500 mt-0.5"
            />
            <div class="ml-2">
                <span class="text-sm text-gray-700 dark:text-gray-300">{props.label}</span>
                {props.description && (
                    <p class="text-xs text-gray-500 dark:text-gray-400">{props.description}</p>
                )}
            </div>
        </label>
    );
}

export interface SelectFieldProps {
    label: string;
    value: string;
    onChange: (value: string) => void;
    options: { value: string; label: string }[];
    description?: string;
    ariaLabel?: string;
}

export function SelectField(props: SelectFieldProps) {
    return (
        <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                {props.label}
            </label>
            <select
                aria-label={props.ariaLabel || props.label}
                value={props.value}
                onChange={(e) => props.onChange(e.currentTarget.value)}
                class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-emerald-500 focus:border-emerald-500"
            >
                {props.options.map(option => (
                    <option value={option.value}>{option.label}</option>
                ))}
            </select>
            {props.description && (
                <p class="text-xs text-gray-500 dark:text-gray-400 mt-1">
                    {props.description}
                </p>
            )}
        </div>
    );
}

export interface ConfigSectionProps {
    title: string;
    children: JSX.Element;
}

export function ConfigSection(props: ConfigSectionProps) {
    return (
        <div class="space-y-4">
            <h5 class="text-md font-medium text-gray-900 dark:text-white">{props.title}</h5>
            {props.children}
        </div>
    );
}

interface ServiceInstanceConfigModalProps {
    isOpen: () => boolean;
    onClose: () => void;
    onSave: (instanceId: string, instanceName: string, syncSettings: ExtendedServiceInstanceSyncSettings) => void;
    serviceInstance: () => ServiceInstance | null;
    customConfigSection?: JSX.Element;
    getServiceSpecificConfig?: () => ServiceSpecificConfig;
    initializeServiceSpecificConfig?: (config: ServiceSpecificConfig | undefined) => void;
    resetServiceSpecificConfig?: () => void;
}

export default function ServiceInstanceConfigModal(props: ServiceInstanceConfigModalProps) {
    const [instanceName, setInstanceName] = createSignal<string>("");
    const [syncFrequency, setSyncFrequency] = createSignal<string>("Every 15 minutes");
    const [conflictResolution, setConflictResolution] = createSignal<string>("Keep both");

    const initializeForm = () => {
        const instance = props.serviceInstance();
        if (instance) {
            setInstanceName(instance.instanceName || instance.name);
            setSyncFrequency(instance.syncSettings?.frequency || "Every 15 minutes");
            setConflictResolution(instance.syncSettings?.conflictResolution || "Keep both");

            // Initialize service-specific settings if handler provided
            const extended = instance.syncSettings as ExtendedServiceInstanceSyncSettings;
            props.initializeServiceSpecificConfig?.(extended?.serviceSpecific);
        }
    };

    const resetForm = () => {
        setInstanceName("");
        setSyncFrequency("Every 15 minutes");
        setConflictResolution("Keep both");
        props.resetServiceSpecificConfig?.();
    };

    const handleClose = () => {
        resetForm();
        props.onClose();
    };

    const handleSave = () => {
        const instance = props.serviceInstance();
        if (!instance) return;

        const serviceSpecific = props.getServiceSpecificConfig?.() || {};

        const syncSettings: ExtendedServiceInstanceSyncSettings = {
            frequency: syncFrequency(),
            conflictResolution: conflictResolution(),
            serviceSpecific
        };

        props.onSave(instance.instanceId, instanceName(), syncSettings);
        props.onClose();
    };

    // Initialize form when modal opens
    const isOpen = () => {
        if (props.isOpen()) {
            initializeForm();
        }
        return props.isOpen();
    };

    const instance = () => props.serviceInstance();

    return (
        <Modal
            isOpen={isOpen}
            onClose={handleClose}
            title={`Configure ${instance()?.instanceName || instance()?.name || 'Service'}`}
            size="lg"
        >
            <div class="space-y-6">
                {/* Instance Name */}
                <div>
                    <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                        Instance Name
                    </label>
                    <input
                        type="text"
                        value={instanceName()}
                        onInput={(e) => setInstanceName(e.currentTarget.value)}
                        placeholder={instance()?.name || "Service name"}
                        class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-white placeholder-gray-500 dark:placeholder-gray-400 focus:ring-2 focus:ring-emerald-500 focus:border-emerald-500"
                    />
                    <p class="text-xs text-gray-500 dark:text-gray-400 mt-1">
                        Give this service instance a custom name to help identify it
                    </p>
                </div>

                {/* Basic Sync Settings */}
                <div class="space-y-4">
                    <h4 class="text-lg font-medium text-gray-900 dark:text-white">Basic Sync Settings</h4>

                    {/* Sync Frequency */}
                    <SelectField
                        label="Sync Frequency"
                        value={syncFrequency()}
                        onChange={setSyncFrequency}
                        ariaLabel="Sync Frequency"
                        options={[
                            { value: "Every 5 minutes", label: "Every 5 minutes" },
                            { value: "Every 15 minutes", label: "Every 15 minutes" },
                            { value: "Every 30 minutes", label: "Every 30 minutes" },
                            { value: "Every hour", label: "Every hour" },
                            { value: "Every 6 hours", label: "Every 6 hours" },
                            { value: "Daily", label: "Daily" },
                            { value: "Manual only", label: "Manual only" }
                        ]}
                    />

                    {/* Conflict Resolution */}
                    <SelectField
                        label="Conflict Resolution"
                        value={conflictResolution()}
                        onChange={setConflictResolution}
                        ariaLabel="Conflict Resolution"
                        description="How to handle conflicts when the same item exists in multiple sources"
                        options={[
                            { value: "Keep both", label: "Keep both" },
                            { value: "Merge when possible", label: "Merge when possible" }
                        ]}
                    />
                </div>

                {/* Service-Specific Configuration Section */}
                <Show when={props.customConfigSection}>
                    <div class="border-t border-gray-200 dark:border-gray-700 pt-6">
                        {props.customConfigSection}
                    </div>
                </Show>

                {/* Current Sync Status */}
                <Show when={instance()?.connected}>
                    <div class="bg-gray-50 dark:bg-gray-700/50 rounded-lg p-4">
                        <div class="flex items-center justify-between">
                            <div>
                                <h5 class="text-sm font-medium text-gray-900 dark:text-white">Sync Status</h5>
                                <p class="text-xs text-gray-500 dark:text-gray-400">
                                    {instance()?.lastSync ? `Last synced: ${instance()?.lastSync}` : 'Never synced'}
                                </p>
                            </div>
                            <div class="flex items-center gap-2">
                                <div class="w-2 h-2 bg-green-500 rounded-full"></div>
                                <span class="text-xs text-gray-500 dark:text-gray-400">Active</span>
                            </div>
                        </div>
                    </div>
                </Show>

                {/* Action Buttons */}
                <div class="flex justify-end space-x-3 pt-4 border-t border-gray-200 dark:border-gray-700">
                    <Button variant="secondary" onClick={handleClose}>
                        Cancel
                    </Button>
                    <Button variant="primary" onClick={handleSave}>
                        <Icon name="v-icon" class="w-4 h-4 mr-2" />
                        Save Configuration
                    </Button>
                </div>
            </div>
        </Modal>
    );
}