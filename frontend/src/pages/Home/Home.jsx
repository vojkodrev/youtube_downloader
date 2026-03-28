import { useState } from 'react'
import { Tv, Search } from 'lucide-react'

const TEST_ITEMS = [
    'Big Buck Bunny',
    'Elephant Dream',
    'Sintel',
    'Tears of Steel',
    'Cosmos Laundromat',
]

export default function Home() {
    const [query, setQuery] = useState('')
    const [focused, setFocused] = useState(false)

    const results = query.trim()
        ? TEST_ITEMS.filter(item => item.toLowerCase().includes(query.toLowerCase()))
        : []

    return (
        <div className="flex flex-col h-[125rem]">

            {/* Top */}
            <div className="bg-gray-900 px-6 py-4 flex items-center gap-3">
                <Tv size={24} className="text-white" />
                <span className="text-white text-xl font-bold tracking-wide">Watch</span>
                <div className="flex-1 flex justify-center">
                    <div className="relative w-full max-w-md">
                        <div className="flex items-center bg-gray-600 rounded-full px-4 py-1.5 gap-2 focus-within:ring-2 focus-within:ring-gray-400">
                            <Search size={16} className="text-gray-300" />
                            <input
                                type="text"
                                placeholder="Search..."
                                value={query}
                                onChange={e => setQuery(e.target.value)}
                                onFocus={() => setFocused(true)}
                                onBlur={() => setFocused(false)}
                                className="bg-transparent text-white placeholder-gray-300 text-sm outline-none w-full"
                            />
                        </div>
                        {focused && results.length > 0 && (
                            <ul className="absolute top-full mt-1 w-full bg-white border border-gray-200 rounded-lg overflow-hidden shadow-lg z-10">
                                {results.map(item => (
                                    <li
                                        key={item}
                                        className="px-4 py-2 text-sm text-gray-900 hover:bg-gray-100 cursor-pointer"
                                    >
                                        {item}
                                    </li>
                                ))}
                            </ul>
                        )}
                    </div>
                </div>
            </div>

            {/* Middle */}
            <div className="flex flex-col md:flex-row flex-1">

                {/* Left: video content + info panel */}
                <div className="flex flex-col md:flex-1">
                    <div className="h-[25rem] md:h-[45rem] bg-white p-4">
                        <p>Video content</p>
                    </div>
                    <div className="h-[10rem] bg-gray-100 p-4">
                        <p>Info panel</p>
                    </div>
                </div>

                {/* Sidebar */}
                <div className="h-[31.25rem] md:h-auto md:w-80 bg-gray-50 p-4">
                    <p>Sidebar</p>
                </div>

            </div>
        </div>
    )
}
