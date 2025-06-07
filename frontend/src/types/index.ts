export type IconName =
	| 'sync-icon'
	| 'menu-icon'
	| 'calendar-icon'
	| 'music-icon'
	| 'realtime-sync-icon'
	| 'lightning-icon'
	| 'lock-icon'
	| 'sun-icon'
	| 'moon-icon'
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
