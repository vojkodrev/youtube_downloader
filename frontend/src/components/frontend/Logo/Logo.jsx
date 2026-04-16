import { Tv } from 'lucide-react'

export default function Logo({ color = 'text-white' }) {
    return (
        <div className="flex items-center gap-3 md:gap-2 lg:gap-3">
            <Tv size={24} className={`${color} -mt-1 md:hidden lg:block`} />
            <Tv size={20} className={`${color} -mt-1 hidden md:block lg:hidden`} />
            <span className={`${color} text-xl md:text-base lg:text-xl font-bold tracking-wide`}>Watch</span>
        </div>
    )
}
