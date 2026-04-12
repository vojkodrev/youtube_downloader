import { Tv } from 'lucide-react'

export default function Logo({ color = 'text-white' }) {
    return (
        <div className="flex items-center gap-3">
            <Tv size={24} className={`${color} -mt-1`} />
            <span className={`${color} text-xl font-bold tracking-wide`}>Watch</span>
        </div>
    )
}
