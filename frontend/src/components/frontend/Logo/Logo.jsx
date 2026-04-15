import { Tv } from 'lucide-react'

export default function Logo({ color = 'text-white' }) {
    return (
        <div className="flex items-center gap-2 lg:gap-3">
            <Tv size={20} className={`${color} -mt-1 lg:hidden`} />
            <Tv size={24} className={`${color} -mt-1 hidden lg:block`} />
            <span className={`${color} text-base lg:text-xl font-bold tracking-wide`}>Watch</span>
        </div>
    )
}
