export type IconName =
	| 'sync-icon'
	| 'menu-icon'
	| 'calendar-icon'
	| 'music-icon'
	| 'realtime-sync-icon'
	| 'lightning-icon'
	| 'location-icon'
	| 'lock-icon'
	| 'mail-icon'
	| 'sun-icon'
	| 'moon-icon'
	| 'phone-icon'
	| 'user-icon'
	| 'logout-icon'
	| 'v-icon'
	| 'google-icon'
	| 'arrow-left'
	| 'clock-icon'
	| 'plus-icon'
	| 'trash-icon'
	// Calendar Services
	| 'google-calendar-icon'
	| 'outlook-icon'
	// Music Services
	| 'spotify-icon'
	| 'apple-music-icon'
	| 'youtube-icon'
	| 'deezer-icon'
	| 'tidal-icon';

export type User = {
	id: string;
	username: string;
	email: string;
	profilePicture?: string;
};

export type ServiceInstanceSyncSettings = {
	frequency: string;
	conflictResolution: string;
};
export interface ServiceInstance {
	instanceId: string;
	instanceName?: string;
	name: string;
	connected: boolean;
	lastSync?: string;
	syncSettings?: ServiceInstanceSyncSettings;
}

export type ServiceType = {
	id: string;
	name: string;
	description: string;
	icon: IconName;
};

export interface ServiceSpecificConfig {
	[key: string]: any;
}

export interface ExtendedServiceInstanceSyncSettings
	extends ServiceInstanceSyncSettings {
	serviceSpecific?: ServiceSpecificConfig;
}
