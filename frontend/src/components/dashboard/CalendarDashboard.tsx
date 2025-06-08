import { createSignal } from "solid-js";
import { BaseDashboard } from "~/components/dashboard";
import type { ServiceInstance, ServiceType } from "~/types";
import {
    CheckboxField,
    ConfigSection,
    SelectField,
    ServiceInstanceConfigModal,
    type ExtendedServiceInstanceSyncSettings,
    type ServiceSpecificConfig
} from ".";
import CalendarConfigModal from "../modals/CalendarConfigModal";

type CalendarConfig = {
    serviceType: string;
    instanceName: string;
    syncFrequency: string;
    conflictResolution: string;
}

interface CalendarService extends ServiceInstance {
    id: string;
    description: string;
    icon: "google-calendar-icon" | "outlook-icon";
}

type CalendarDashboardProps = {
    onBack: () => void;
}

export default function CalendarDashboard(props: CalendarDashboardProps) {
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

    const [services, setServices] = createSignal<CalendarService[]>([
        {
            id: "outlook",
            instanceId: "outlook-1",
            name: "Microsoft Outlook",
            instanceName: "Work Account",
            description: "Connect your Outlook calendar and Exchange events",
            icon: "outlook-icon",
            connected: true,
            lastSync: "2 minutes ago",
            syncSettings: {
                frequency: "Every 15 minutes",
                conflictResolution: "Keep both"
            }
        }
    ]);

    const [isAddModalOpen, setIsAddModalOpen] = createSignal(false);
    const [isConfigModalOpen, setIsConfigModalOpen] = createSignal(false);
    const [selectedServiceForConfig, setSelectedServiceForConfig] = createSignal<CalendarService | null>(null);

    // Calendar-specific configuration signals
    const [syncEvents, setSyncEvents] = createSignal<boolean>(true);
    const [syncTasks, setSyncTasks] = createSignal<boolean>(false);
    const [defaultCalendar, setDefaultCalendar] = createSignal<string>("primary");
    const [reminderSettings, setReminderSettings] = createSignal<boolean>(true);
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

    const addServiceInstance = (serviceTypeId: string) => {
        const serviceType = serviceTypes.find(t => t.id === serviceTypeId);
        if (!serviceType) return;

        const existingInstances = services().filter(s => s.id === serviceTypeId);
        const instanceNumber = existingInstances.length + 1;

        const newService: CalendarService = {
            id: serviceType.id,
            name: serviceType.name,
            description: serviceType.description,
            icon: serviceType.icon as CalendarService["icon"],
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

    const handleModalSave = (config: CalendarConfig) => {
        const serviceType = serviceTypes.find(t => t.id === config.serviceType);
        if (!serviceType) return;

        const existingInstances = services().filter(s => s.id === config.serviceType);
        const instanceNumber = existingInstances.length + 1;

        const newService: CalendarService = {
            id: serviceType.id,
            name: serviceType.name,
            description: serviceType.description,
            icon: serviceType.icon as CalendarService["icon"],
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

    const handleSync = (instanceId: string) => {
        console.log("Syncing instance", instanceId);
    };

    const removeServiceInstance = (instanceId: string) => {
        setServices(prev => prev.filter(service => service.instanceId !== instanceId));
    };

    const getServicesByType = (typeId: string) => {
        return services().filter(service => service.id === typeId);
    };

    const connectedCount = () => services().filter(s => s.connected).length;

    // Calendar-specific configuration handlers
    const getServiceSpecificConfig = (): ServiceSpecificConfig => {
        return {
            syncEvents: syncEvents(),
            syncTasks: syncTasks(),
            defaultCalendar: defaultCalendar(),
            reminderSettings: reminderSettings(),
            bidirectionalSync: bidirectionalSync(),
            autoResolveConflicts: autoResolveConflicts(),
            privateMode: privateMode()
        };
    };

    const initializeServiceSpecificConfig = (config: ServiceSpecificConfig | undefined) => {
        setSyncEvents(config?.syncEvents ?? true);
        setSyncTasks(config?.syncTasks ?? false);
        setDefaultCalendar(config?.defaultCalendar || "primary");
        setReminderSettings(config?.reminderSettings ?? true);
        setBidirectionalSync(config?.bidirectionalSync ?? true);
        setAutoResolveConflicts(config?.autoResolveConflicts ?? false);
        setPrivateMode(config?.privateMode ?? false);
    };

    const resetServiceSpecificConfig = () => {
        setSyncEvents(true);
        setSyncTasks(false);
        setDefaultCalendar("primary");
        setReminderSettings(true);
        setBidirectionalSync(true);
        setAutoResolveConflicts(false);
        setPrivateMode(false);
    };

    // Calendar-specific configuration section
    const calendarConfigSection = (
        <>
            <ConfigSection title="Calendar Sync Options">
                <div class="space-y-3">
                    <CheckboxField
                        label="Sync Events"
                        checked={syncEvents()}
                        onChange={setSyncEvents}
                    />
                    <CheckboxField
                        label="Sync Tasks & To-dos"
                        checked={syncTasks()}
                        onChange={setSyncTasks}
                    />
                    <CheckboxField
                        label="Sync Reminder Settings"
                        checked={reminderSettings()}
                        onChange={setReminderSettings}
                    />
                </div>

                <SelectField
                    label="Default Calendar"
                    value={defaultCalendar()}
                    onChange={setDefaultCalendar}
                    ariaLabel="Default Calendar"
                    options={[
                        { value: "primary", label: "Primary Calendar" },
                        { value: "work", label: "Work Calendar" },
                        { value: "personal", label: "Personal Calendar" },
                        { value: "shared", label: "Shared Calendar" }
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
            title="Calendar Integration"
            description="Connect and manage your calendar services"
            addButtonText="Add Calendar"
            mainIcon="calendar-icon"
            serviceTypes={serviceTypes}
            services={services}
            connectedCount={connectedCount}
            lastSyncValue="2m"
            onBack={props.onBack}
            onAdd={() => setIsAddModalOpen(true)}
            onAddAccount={addServiceInstance}
            onSync={handleSync}
            onConnect={handleConnect}
            onRemove={removeServiceInstance}
            onConfigure={handleConfigure}
            getServicesByType={getServicesByType}
            modal={
                <>
                    <CalendarConfigModal
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
                        customConfigSection={calendarConfigSection}
                        getServiceSpecificConfig={getServiceSpecificConfig}
                        initializeServiceSpecificConfig={initializeServiceSpecificConfig}
                        resetServiceSpecificConfig={resetServiceSpecificConfig}
                    />
                </>
            }
        />
    );
} 