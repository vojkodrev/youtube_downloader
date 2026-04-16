import { DownloadIcon } from 'lucide-react'
import Logo from '@/components/frontend/Logo/Logo'
import SidebarTrigger from '@/components/frontend/SidebarTrigger/SidebarTrigger'
import {
    Sidebar,
    SidebarHeader,
    SidebarContent,
    SidebarGroup,
    SidebarGroupContent,
    SidebarMenu,
    SidebarMenuItem,
    SidebarMenuButton,
} from '@/components/ui/sidebar'

export default function AppSidebar({ onRequestDownload }) {
    return (
        <Sidebar collapsible="offcanvas">
            <SidebarHeader className="flex flex-row items-center gap-3 px-3 py-4">
                <SidebarTrigger className="text-black md:ml-3" />
                <Logo color="text-black" />
            </SidebarHeader>
            <SidebarContent>
                <SidebarGroup>
                    <SidebarGroupContent>
                        <SidebarMenu>
                            <SidebarMenuItem>
                                <SidebarMenuButton onClick={onRequestDownload}>
                                    <DownloadIcon />
                                    <span>Request Download</span>
                                </SidebarMenuButton>
                            </SidebarMenuItem>
                        </SidebarMenu>
                    </SidebarGroupContent>
                </SidebarGroup>
            </SidebarContent>
        </Sidebar>
    )
}
