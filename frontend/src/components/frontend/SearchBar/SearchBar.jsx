import { useState } from 'react'
import { Search, X } from 'lucide-react'

const TEST_ITEMS = [
    'Big Buck Bunny',
    'Elephant Dream',
    'Sintel',
    'Tears of Steel',
    'Cosmos Laundromat',
]

export default function SearchBar() {
    const [query, setQuery] = useState('')
    const [focused, setFocused] = useState(false)

    const results = query.trim()
        ? TEST_ITEMS.filter(item => item.toLowerCase().includes(query.toLowerCase()))
        : []

    return (
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
                    {query && (
                        <button onMouseDown={e => { e.preventDefault(); setQuery('') }}>
                            <X size={14} className="text-gray-300 hover:text-white" />
                        </button>
                    )}
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
    )
}
