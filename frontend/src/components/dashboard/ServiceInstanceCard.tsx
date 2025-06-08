import Button from "~/components/ui/Button";
import Icon from "~/components/ui/Icon";
import type { ServiceInstance } from "~/types";



interface ServiceInstanceCardProps {
    service: ServiceInstance;
    onConnect: (instanceId: string) => void;
    onRemove: (instanceId: string) => void;
    onConfigure?: (instanceId: string) => void;
    onSync?: (instanceId: string) => void;
}

export default function ServiceInstanceCard(props: ServiceInstanceCardProps) {
    return (
        <div class="bg-white dark:bg-gray-800 rounded-xl border border-gray-200 dark:border-gray-700 p-6">
            <div class="flex items-start justify-between mb-4">
                <div>
                    <h4 class="text-base font-medium text-gray-900 dark:text-white">
                        {props.service.instanceName || props.service.name}
                    </h4>
                    {props.service.instanceName && (
                        <p class="text-sm text-gray-500 dark:text-gray-400">
                            {props.service.name}
                        </p>
                    )}
                </div>
                <div class="flex items-center gap-2">
                    <div class={`px-2 py-1 rounded-full text-xs font-medium ${props.service.connected
                        ? 'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400'
                        : 'bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-300'
                        }`}>
                        {props.service.connected ? 'Connected' : 'Disconnected'}
                    </div>
                    <Button
                        variant="secondary"
                        size="sm"
                        onClick={() => props.onRemove(props.service.instanceId)}
                        class="p-1"
                    >
                        <Icon name="trash-icon" class="w-4 h-4" />
                    </Button>
                </div>
            </div>

            {props.service.lastSync && (
                <p class="text-xs text-gray-500 dark:text-gray-400 mb-4">
                    Last synced: {props.service.lastSync}
                </p>
            )}

            <div class="flex gap-3">
                <Button
                    variant={props.service.connected ? "secondary" : "primary"}
                    size="sm"
                    onClick={() => props.onConnect(props.service.instanceId)}
                    class="flex-1"
                >
                    {props.service.connected ? 'Disconnect' : 'Connect'}
                </Button>
                {props.service.connected && props.onSync && (
                    <Button
                        variant="secondary"
                        size="sm"
                        onClick={() => props.onSync!(props.service.instanceId)}
                    >
                        <Icon name="sync-icon" class="w-4 h-4 mr-1" />
                        Sync
                    </Button>
                )}
                {props.service.connected && props.onConfigure && (
                    <Button
                        variant="secondary"
                        size="sm"
                        onClick={() => props.onConfigure!(props.service.instanceId)}
                    >
                        Configure
                    </Button>
                )}
            </div>
        </div>
    );
} 