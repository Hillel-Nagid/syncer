import { JSX } from "solid-js";
import type { IconName } from "~/types";
import {
    AppleMusicIcon,
    ArrowLeftIcon,
    CalendarIcon,
    ClockIcon,
    DeezerIcon,
    // Brand icons
    GoogleCalendarIcon,
    GoogleIcon,
    LightningIcon,
    LockIcon,
    LogoutIcon,
    MenuIcon,
    MoonIcon,
    MusicIcon,
    OutlookIcon,
    PlusIcon,
    RealtimeSyncIcon,
    // Music Services
    SpotifyIcon,
    SunIcon,
    SyncIcon,
    TidalIcon,
    TrashIcon,
    UserIcon,
    VIcon,
    YouTubeIcon
} from "../icons";

interface IconProps {
    name: IconName;
    class?: string;
    alt?: string;
}

const iconMap: Record<IconName, (props: { class?: string }) => JSX.Element> = {
    "sync-icon": SyncIcon,
    "menu-icon": MenuIcon,
    "calendar-icon": CalendarIcon,
    "music-icon": MusicIcon,
    "realtime-sync-icon": RealtimeSyncIcon,
    "lightning-icon": LightningIcon,
    "lock-icon": LockIcon,
    "sun-icon": SunIcon,
    "moon-icon": MoonIcon,
    "user-icon": UserIcon,
    "logout-icon": LogoutIcon,
    "v-icon": VIcon,
    "google-icon": GoogleIcon,
    "arrow-left": ArrowLeftIcon,
    "clock-icon": ClockIcon,
    "plus-icon": PlusIcon,
    "trash-icon": TrashIcon,
    // Calendar Services
    "google-calendar-icon": GoogleCalendarIcon,
    "outlook-icon": OutlookIcon,
    // Music Services
    "spotify-icon": SpotifyIcon,
    "apple-music-icon": AppleMusicIcon,
    "youtube-icon": YouTubeIcon,
    "deezer-icon": DeezerIcon,
    "tidal-icon": TidalIcon,
};

export default function Icon(props: IconProps): JSX.Element {
    const IconComponent = iconMap[props.name];

    if (!IconComponent) {
        console.warn(`Icon "${props.name}" not found`);
        return <span>?</span>;
    }

    return <IconComponent class={props.class} />;
} 