import { useEffect } from 'react'
import { MenuIcon } from 'lucide-react'
import { useSidebar } from '@/components/ui/sidebar'

export default function SidebarTrigger({ className = 'text-white' }) {
    const { toggleSidebar, open, setOpen } = useSidebar()

    useEffect(() => {
        function handleKeyDown(e) {
            if (e.key === 'Escape' && open) setOpen(false)
        }
        document.addEventListener('keydown', handleKeyDown)
        return () => document.removeEventListener('keydown', handleKeyDown)
    }, [open, setOpen])

    return (
        <button onClick={toggleSidebar} className={`lg:pr-3 ${className}`}>
            <MenuIcon size={20} className="lg:hidden" />
            <MenuIcon size={24} className="hidden lg:block" />
        </button>
    )
}
