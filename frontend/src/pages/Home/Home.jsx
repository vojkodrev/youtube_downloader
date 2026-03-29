import { useState, useEffect } from 'react'
import { useParams, Link, useNavigate, useSearchParams } from 'react-router-dom'
import Logo from '@/components/frontend/Logo/Logo'
import SearchBar from '@/components/frontend/SearchBar/SearchBar'

const API_URL = import.meta.env.VITE_API_URL

export default function Home() {
    const { id } = useParams()
    const [searchParams] = useSearchParams()
    const navigate = useNavigate()
    const [videos, setVideos] = useState([])
    const [selectedVideo, setSelectedVideo] = useState(null)
    const [duration, setDuration] = useState(null)

    useEffect(() => {
        (async () => {
            const res = await fetch(`${API_URL}/videos`)
            const data = await res.json()
            const dataWithTimes = data.map(v => ({
                ...v,
                savedTime: localStorage.getItem(`time_${v.id}`)
            }))
            setVideos(dataWithTimes)
        })()
    }, [])

    useEffect(() => {
        if (videos.length === 0) return
        if (!id) {
            navigate(`/watch/${videos[0].id}`, { replace: true })
            return
        }
        if (!searchParams.get('t')) {
            const savedTime = Math.floor(parseFloat(localStorage.getItem(`time_${id}`)))
            if (savedTime && savedTime > 10) {
                navigate(`/watch/${id}?t=${savedTime}`, { replace: true })
                return
            }
        }
        const match = videos.find(v => v.id === id)
        setSelectedVideo(match)
        setDuration(null)
    }, [id, videos])

    return (
        <div className="flex flex-col">

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
                                onTimeUpdate={e => {
                                    const t = e.target.currentTime
                                    if (Math.floor(t) % 5 !== 0) return
                                    localStorage.setItem(`time_${selectedVideo.id}`, Math.floor(t))
                                    setVideos(prev => prev.map(v => v.id === selectedVideo.id ? { ...v, savedTime: t } : v))
                                }}
                                onLoadedMetadata={e => {
                                    setDuration(e.target.duration)
                                    const t = searchParams.get('t')
                                    if (t) e.target.currentTime = parseFloat(t)
                                }}
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
                        <Link
                            key={video.id}
                            to={video.savedTime ? `/watch/${video.id}?t=${video.savedTime}` : `/watch/${video.id}`}
                            title={video.name}
                            className={`flex gap-2 p-2 hover:bg-gray-100 ${selectedVideo?.id === video.id ? 'bg-gray-200' : ''}`}
                        >
                            <img
                                src={`${API_URL}/thumbnail/${video.id}`}
                                className="w-36 h-20 object-cover rounded flex-shrink-0 bg-gray-300"
                            />
                            <div className="flex flex-col justify-center min-w-0">
                                <p className="text-sm font-medium truncate text-gray-900">
                                    {video.name}
                                </p>
                                <p className="text-xs text-gray-500 mt-1">
                                    {new Date(video.date).toLocaleString()}
                                </p>
                            </div>
                        </Link>
                    ))}
                </div>

            </div>
        </div>
    )
}
