import { useState } from 'react'
import { ChevronDownIcon, ClockIcon, DownloadIcon, LogOutIcon } from 'lucide-react'
import { useNavigate } from 'react-router-dom'
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

export default function AppSidebar({ onRequestVideoDownload, onRequestPlaylistDownload, onExportWatchedTimes, onImportWatchedTimes }) {
    const [downloadOpen, setDownloadOpen] = useState(false)
    const [watchedTimesOpen, setWatchedTimesOpen] = useState(false)
    const navigate = useNavigate()

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
                                <SidebarMenuButton onClick={() => setDownloadOpen(o => !o)}>
                                    <DownloadIcon />
                                    <span>Request Download</span>
                                    <ChevronDownIcon className={`ml-auto transition-transform ${downloadOpen ? 'rotate-180' : ''}`} />
                                </SidebarMenuButton>
                                {downloadOpen && (
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
                            <SidebarMenuItem>
                                <SidebarMenuButton onClick={() => setWatchedTimesOpen(o => !o)}>
                                    <ClockIcon />
                                    <span>Watched Times</span>
                                    <ChevronDownIcon className={`ml-auto transition-transform ${watchedTimesOpen ? 'rotate-180' : ''}`} />
                                </SidebarMenuButton>
                                {watchedTimesOpen && (
                                    <SidebarMenuSub>
                                        <SidebarMenuSubItem>
                                            <SidebarMenuSubButton onClick={() => onExportWatchedTimes()}>
                                                <span>Export</span>
                                            </SidebarMenuSubButton>
                                        </SidebarMenuSubItem>
                                        <SidebarMenuSubItem>
                                            <SidebarMenuSubButton onClick={() => onImportWatchedTimes()}>
                                                <span>Import</span>
                                            </SidebarMenuSubButton>
                                        </SidebarMenuSubItem>
                                    </SidebarMenuSub>
                                )}
                            </SidebarMenuItem>
                        </SidebarMenu>
                    </SidebarGroupContent>
                </SidebarGroup>
                <SidebarGroup>
                    <SidebarGroupContent>
                        <SidebarMenu>
                            <SidebarMenuItem>
                                <SidebarMenuButton onClick={() => navigate('/logout')}>
                                    <LogOutIcon />
                                    <span>Logout</span>
                                </SidebarMenuButton>
                            </SidebarMenuItem>
                        </SidebarMenu>
                    </SidebarGroupContent>
                </SidebarGroup>
            </SidebarContent>
        </Sidebar>
    )
}
