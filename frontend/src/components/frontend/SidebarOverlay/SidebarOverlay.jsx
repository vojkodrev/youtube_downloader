import { useSidebar } from '@/components/ui/sidebar'

export default function SidebarOverlay() {
    const { open, setOpen } = useSidebar()
    if (!open) return null
    return <div className="fixed inset-0 z-9 bg-black/50" onClick={() => setOpen(false)} />
}
