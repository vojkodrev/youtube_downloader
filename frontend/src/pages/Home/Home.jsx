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
                <div className="h-[31.25rem] md:h-auto md:w-80 bg-gray-50 p-4">
                    <p>Sidebar</p>
                </div>

            </div>
        </div>
    )
}
