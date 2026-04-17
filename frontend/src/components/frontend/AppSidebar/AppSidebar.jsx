import { useState } from 'react'
import { ChevronDownIcon, DownloadIcon } from 'lucide-react'
import Logo from '@/components/frontend/Logo/Logo'
import SidebarTrigger from '@/components/frontend/SidebarTrigger/SidebarTrigger'
import {
    Sidebar,
    SidebarHeader,
    SidebarContent,
    SidebarGroup,
    SidebarGroupContent,
    SidebarGroupLabel,
    SidebarMenu,
    SidebarMenuItem,
    SidebarMenuButton,
    SidebarMenuSub,
    SidebarMenuSubItem,
    SidebarMenuSubButton,
} from '@/components/ui/sidebar'

export default function AppSidebar({ onRequestVideoDownload, onRequestPlaylistDownload }) {
    const [open, setOpen] = useState(false)

    return (
        <Sidebar collapsible="offcanvas">
            <SidebarHeader className="flex flex-row items-center gap-3 px-3 py-4">
                <SidebarTrigger className="text-black md:ml-3" />
                <Logo color="text-black" />
            </SidebarHeader>
            <SidebarContent>
                <SidebarGroup>
                    <SidebarGroupLabel>Videos</SidebarGroupLabel>
                    <SidebarGroupContent>
                        <SidebarMenu>
                            <SidebarMenuItem>
                                <SidebarMenuButton onClick={() => setOpen(o => !o)}>
                                    <DownloadIcon />
                                    <span>Request Download</span>
                                    <ChevronDownIcon className={`ml-auto transition-transform ${open ? 'rotate-180' : ''}`} />
                                </SidebarMenuButton>
                                {open && (
                                    <SidebarMenuSub>
                                        <SidebarMenuSubItem>
                                            <SidebarMenuSubButton onClick={() => onRequestVideoDownload()}>
                                                <span>Video</span>
                                            </SidebarMenuSubButton>
                                        </SidebarMenuSubItem>
                                        <SidebarMenuSubItem>
                                            <SidebarMenuSubButton onClick={() => onRequestPlaylistDownload()}>
                                                <span>Playlist</span>
                                            </SidebarMenuSubButton>
                                        </SidebarMenuSubItem>
                                    </SidebarMenuSub>
                                )}
                            </SidebarMenuItem>
                        </SidebarMenu>
                    </SidebarGroupContent>
                </SidebarGroup>
            </SidebarContent>
        </Sidebar>
    )
}
