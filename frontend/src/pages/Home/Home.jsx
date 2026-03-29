import { useState, useEffect } from 'react'
import Logo from '@/components/frontend/Logo/Logo'
import SearchBar from '@/components/frontend/SearchBar/SearchBar'

const API_URL = import.meta.env.VITE_API_URL

export default function Home() {
    const [videos, setVideos] = useState([])
    const [selectedVideo, setSelectedVideo] = useState(null)
    const [duration, setDuration] = useState(null)

    useEffect(() => {
        (async () => {
            const res = await fetch(`${API_URL}/videos`)
            const data = await res.json()
            setVideos(data)
            if (data.length > 0) setSelectedVideo(data[0])
        })()
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
                    <div className="h-[25rem] md:h-[45rem] bg-black">
                        {selectedVideo && (
                            <video
                                key={selectedVideo.id}
                                src={`${API_URL}/video/${selectedVideo.id}`}
                                controls
                                autoPlay
                                onLoadedMetadata={e => setDuration(e.target.duration)}
                                className="w-full h-full"
                            />
                        )}
                    </div>
                    <div className="bg-gray-100 p-4">
                        {selectedVideo && (
                            <>
                                <p className="font-semibold text-lg truncate" title={selectedVideo.name}>{selectedVideo.name}</p>
                                {duration != null && (
                                    <p className="text-sm text-gray-500 mt-1">
                                        {new Date(duration * 1000).toISOString().substring(11, 19)}
                                    </p>
                                )}
                                <p className="text-sm text-gray-500 mt-1">{new Date(selectedVideo.date).toLocaleString()}</p>
                            </>
                        )}
                    </div>
                </div>

                {/* Sidebar */}
                <div className="h-[31.25rem] md:h-auto md:w-[28rem] bg-gray-50 overflow-y-auto">
                    {videos.map(video => (
                        <div
                            key={video.id}
                            title={video.name}
                            onClick={() => { setSelectedVideo(video); setDuration(null) }}
                            className={`flex gap-2 p-2 cursor-pointer hover:bg-gray-100 ${selectedVideo?.id === video.id ? 'bg-gray-200' : ''}`}
                        >
                            <img
                                src={`${API_URL}/thumbnail/${video.id}`}
                                className="w-36 h-20 object-cover rounded flex-shrink-0 bg-gray-300"
                            />
                            <div className="flex flex-col justify-center min-w-0">
                                <p className="text-sm font-medium truncate">
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
