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
    primary_email: string;
    full_name: string;
    avatar_url?: string;
    last_login?: string;
    online?: boolean;
    created_at: string;
    updated_at: string;
};

// Authentication types
export interface LoginRequest {
    email: string;
    password: string;
}

export interface RegisterRequest {
    email: string;
    password: string;
    full_name: string;
}

export interface AuthResponse {
    user: User;
}

export interface VerifyEmailRequest {
    token: string;
}

export interface ResendVerificationRequest {
    email: string;
}

export interface ForgotPasswordRequest {
    email: string;
}

export interface ResetPasswordRequest {
    token: string;
    password: string;
}

export interface RefreshTokenRequest {
    refresh_token: string;
}

export interface RefreshTokenResponse {
    access_token: string;
    refresh_token: string;
    expires_in: number;
}

export interface ApiError {
    error: string;
}

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
