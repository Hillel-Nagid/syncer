import { JSX } from "solid-js";
import type { IconName } from "~/types";
import {
    CalendarIcon,
    MenuIcon,
    MoonIcon,
    MusicIcon,
    RealtimeSyncIcon,
    SunIcon,
    SyncIcon
} from "./icons";

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
    "sun-icon": SunIcon,
    "moon-icon": MoonIcon,
};

export default function Icon(props: IconProps): JSX.Element {
    const IconComponent = iconMap[props.name];

    if (!IconComponent) {
        console.warn(`Icon "${props.name}" not found`);
        return <span>?</span>;
    }

    return <IconComponent class={props.class} />;
} 