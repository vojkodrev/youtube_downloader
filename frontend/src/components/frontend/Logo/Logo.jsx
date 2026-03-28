import { Tv } from 'lucide-react'

export default function Logo() {
    return (
        <div className="flex items-center gap-3">
            <Tv size={24} className="text-white" />
            <span className="text-white text-xl font-bold tracking-wide">Watch</span>
        </div>
    )
}
