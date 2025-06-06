export type IconName =
	| 'sync-icon'
	| 'menu-icon'
	| 'calendar-icon'
	| 'music-icon'
	| 'realtime-sync-icon'
	| 'sun-icon'
	| 'moon-icon'
	| 'user-icon'
	| 'logout-icon'
	| 'v-icon'
	| 'google-icon';

export type User = {
	id: string;
	username: string;
	email: string;
	profilePicture?: string;
};
