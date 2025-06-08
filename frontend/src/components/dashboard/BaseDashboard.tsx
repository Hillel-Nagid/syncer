import { JSX } from "solid-js";
import type { IconName, ServiceInstance, ServiceType } from "~/types";
import DashboardHeader from "./DashboardHeader";
import ServiceTypeSection from "./ServiceTypeSection";
import StatusOverview from "./StatusOverview";

interface BaseDashboardProps<T extends ServiceInstance> {
    title: string;
    description: string;
    addButtonText: string;
    mainIcon: IconName;
    serviceTypes: ServiceType[];
    services: () => T[];
    connectedCount: () => number;
    lastSyncValue: string;
    onBack: () => void;
    onAdd: () => void;
    onAddAccount: (serviceTypeId: string) => void;
    onConnect: (instanceId: string) => void;
    onSync: (instanceId: string) => void;
    onRemove: (instanceId: string) => void;
    onConfigure?: (instanceId: string) => void;
    getServicesByType: (typeId: string) => T[];
    modal: JSX.Element;
}

export default function BaseDashboard<T extends ServiceInstance>(props: BaseDashboardProps<T>) {
    return (
        <div class="min-h-screen w-full bg-gray-50 dark:bg-gray-900">
            <div class="container mx-auto px-6 py-16 max-w-5xl">
                {/* Header */}
                <DashboardHeader
                    title={props.title}
                    description={props.description}
                    addButtonText={props.addButtonText}
                    onBack={props.onBack}
                    onAdd={props.onAdd}
                />

                {/* Status Overview */}
                <StatusOverview
                    connectedCount={props.connectedCount()}
                    mainIcon={props.mainIcon}
                    lastSyncValue={props.lastSyncValue}
                />

                {/* Services */}
                <div class="mb-12">
                    <h2 class="text-2xl font-bold text-gray-900 dark:text-white mb-6">
                        {props.title.split(' ')[0]} Services
                    </h2>
                    <div class="space-y-8">
                        {props.serviceTypes.map((serviceType) => {
                            const instances = props.getServicesByType(serviceType.id);
                            return (
                                <ServiceTypeSection
                                    serviceType={serviceType}
                                    instances={instances}
                                    onConnect={props.onConnect}
                                    onRemove={props.onRemove}
                                    onSync={props.onSync}
                                    onConfigure={props.onConfigure}
                                />
                            );
                        })}
                    </div>
                </div>
            </div>

            {/* Modal */}
            {props.modal}
        </div>
    );
} 