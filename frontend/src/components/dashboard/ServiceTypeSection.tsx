import Icon from "~/components/ui/Icon";
import type { IconName } from "~/types";
import ServiceInstanceCard, { type ServiceInstance } from "./ServiceInstanceCard";

export interface ServiceType {
    id: string;
    name: string;
    description: string;
    icon: IconName;
}

interface ServiceTypeSectionProps {
    serviceType: ServiceType;
    instances: ServiceInstance[];
    onConnect: (instanceId: string) => void;
    onRemove: (instanceId: string) => void;
    onSync?: (instanceId: string) => void;
    onConfigure?: (instanceId: string) => void;
}

export default function ServiceTypeSection(props: ServiceTypeSectionProps) {
    return (
        <div class="space-y-4">
            {/* Service Type Header */}
            <div class="flex items-center justify-between">
                <div class="flex items-center">
                    <div class="w-10 h-10 rounded-lg bg-emerald-100 dark:bg-emerald-900/30 flex items-center justify-center mr-3">
                        <Icon
                            name={props.serviceType.icon}
                            class="w-5 h-5 text-emerald-600 dark:text-emerald-400"
                        />
                    </div>
                    <div>
                        <h3 class="text-lg font-semibold text-gray-900 dark:text-white">
                            {props.serviceType.name}
                        </h3>
                        <p class="text-sm text-gray-600 dark:text-gray-300">
                            {props.serviceType.description}
                        </p>
                    </div>
                </div>
            </div>

            {/* Service Instances */}
            {props.instances.length > 0 ? (
                <div class="grid md:grid-cols-2 gap-4 ml-13">
                    {props.instances.map((service) => (
                        <ServiceInstanceCard
                            service={service}
                            onSync={props.onSync}
                            onConnect={props.onConnect}
                            onRemove={props.onRemove}
                            onConfigure={props.onConfigure}
                        />
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
} 