import Icon from "~/components/ui/Icon";
import type { IconName } from "~/types";

interface StatusOverviewProps {
    connectedCount: number;
    mainIcon: IconName;
    lastSyncValue: string;
}

export default function StatusOverview(props: StatusOverviewProps) {
    return (
        <div class="grid md:grid-cols-3 gap-6 mb-12">
            <div class="bg-white dark:bg-gray-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700">
                <div class="flex items-center justify-between">
                    <div>
                        <p class="text-sm font-medium text-gray-600 dark:text-gray-400">Connected Services</p>
                        <p class="text-2xl font-bold text-emerald-600 dark:text-emerald-400">{props.connectedCount}</p>
                    </div>
                    <Icon name={props.mainIcon} class="w-8 h-8 text-emerald-600 dark:text-emerald-400" />
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
                        <p class="text-2xl font-bold text-blue-600 dark:text-blue-400">{props.lastSyncValue}</p>
                    </div>
                    <Icon name="clock-icon" class="w-8 h-8 text-blue-600 dark:text-blue-400" />
                </div>
            </div>
        </div>
    );
} 