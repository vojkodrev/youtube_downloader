import { useState } from 'react'
import { Search, X } from 'lucide-react'
import { useNavigate } from 'react-router-dom'

function fuzzyMatch(str, query) {
    const s = str.toLowerCase()
    const q = query.toLowerCase()
    let si = 0, qi = 0, score = 0, consecutive = 0
    while (si < s.length && qi < q.length) {
        if (s[si] === q[qi]) {
            score += 1 + consecutive
            consecutive++
            qi++
        } else {
            consecutive = 0
        }
        si++
    }
    return qi === q.length ? score : -1
}

export default function SearchBar({ videos = [] }) {
    const [query, setQuery] = useState('')
    const [focused, setFocused] = useState(false)
    const navigate = useNavigate()

    const results = query.trim()
        ? videos
            .map(v => ({ v, score: fuzzyMatch(v.name, query.trim()) }))
            .filter(({ score }) => score >= 0)
            .sort((a, b) => b.score - a.score)
            .map(({ v }) => v)
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
                {focused && query.trim() && (
                    <ul className="absolute top-full mt-1 w-full bg-white border border-gray-200 rounded-lg overflow-hidden shadow-lg z-10">
                        {results.length > 0
                            ? results.map(video => (
                                <li
                                    key={video.id}
                                    onMouseDown={() => { setQuery(''); navigate(`/watch/${video.id}`) }}
                                    className="px-4 py-2 text-sm text-gray-900 hover:bg-gray-100 cursor-pointer flex items-center gap-2 min-w-0"
                                    title={video.name}
                                >
                                    <Search size={14} className="text-gray-400 shrink-0" />
                                    <span className="truncate">{video.name}</span>
                                </li>
                            ))
                            : <li className="px-4 py-2 text-sm text-gray-400 select-none">No results found</li>
                        }
                    </ul>
                )}
            </div>
        </div>
    )
}
