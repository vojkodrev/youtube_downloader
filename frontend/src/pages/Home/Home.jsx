import { useState, useEffect } from 'react'
import Logo from '@/components/frontend/Logo/Logo'
import SearchBar from '@/components/frontend/SearchBar/SearchBar'

const API_URL = import.meta.env.VITE_API_URL

export default function Home() {
    const [videos, setVideos] = useState([])

    useEffect(() => {
        async function loadVideos() {
            const res = await fetch(`${API_URL}/videos`)
            const data = await res.json()
            setVideos(data)
        }
        loadVideos()
    }, [])

    return (
        <div className="flex flex-col h-[125rem]">

            {/* Top */}
            <div className="bg-gray-900 px-6 py-4 flex items-center gap-3">
                <Logo />
                <SearchBar />
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
                <div className="h-[31.25rem] md:h-auto md:w-[26rem] bg-gray-50 overflow-y-auto">
                    {videos.map(video => (
                        <div key={video.id} title={video.name} className="flex gap-2 p-2 hover:bg-gray-100 cursor-pointer">
                            <img
                                src={`${API_URL}/thumbnail/${video.id}`}
                                className="w-24 h-16 object-cover rounded flex-shrink-0 bg-gray-300"
                            />
                            <div className="flex flex-col justify-center min-w-0">
                                <p
                                    className="text-sm font-medium truncate"
                                >
                                    {video.name}
                                </p>
                                <p className="text-xs text-gray-500 mt-1">
                                    {new Date(video.date).toLocaleString()}
                                </p>
                            </div>
                        </div>
                    ))}
                </div>

            </div>
        </div>
    )
}
