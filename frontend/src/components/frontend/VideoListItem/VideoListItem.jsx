import { useState } from 'react'
import { Link } from 'react-router-dom'
import { EllipsisVertical, Download, X } from 'lucide-react'
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'

const API_URL = import.meta.env.VITE_API_URL

function getDownloadInfo(id) {
    const raw = localStorage.getItem(`download_${id}`)
    return raw ? JSON.parse(raw) : null
}

function saveDownloadInfo(id, filename) {
    localStorage.setItem(`download_${id}`, JSON.stringify({ filename }))
}

export default function VideoListItem({ video, isSelected, infoVisible }) {
    const [downloadInfo, setDownloadInfo] = useState(getDownloadInfo(video.id))

    return (
        <div className={`flex items-center p-2 hover:bg-gray-100 ${isSelected ? 'bg-gray-200' : ''}`}>
            <Link
                to={video.savedTime ? `/watch/${video.id}?t=${video.savedTime}` : `/watch/${video.id}`}
                title={video.name}
                className="flex gap-2 min-w-0 flex-1"
            >
                <div className="flex flex-row gap-2 min-w-0 overflow-hidden">
                    <div className="relative flex-shrink-0">
                        <img
                            src={`${API_URL}/thumbnail/${video.id}`}
                            className="w-36 h-20 object-cover rounded bg-gray-300"
                        />
                        {!!+video.savedTime && video.duration && (
                            <div className="absolute bottom-0 left-0 right-0 h-1 bg-white/30 rounded-bl">
                                <div
                                    className={`h-full bg-red-500 rounded-bl${video.savedTime >= video.duration ? ' rounded-br' : ''}`}
                                    style={{ width: `${Math.min(video.savedTime / video.duration * 100, 100)}%` }}
                                />
                            </div>
                        )}
                        {video.duration != null && (
                            <span className="absolute bottom-2 right-1 bg-black/60 text-white text-xs px-1 rounded">
                                {new Date(video.duration * 1000).toISOString().substring(video.duration >= 3600 ? 11 : 14, 19)}
                            </span>
                        )}
                    </div>
                    <div className="flex flex-col justify-start min-w-0">
                        <p className="text-sm font-medium text-gray-900 line-clamp-2">
                            {video.name}
                        </p>
                        <p className="text-xs text-gray-500 mt-1">
                            {[video.channel, new Date(video.date).toLocaleString()].filter(Boolean).join(' · ')}
                        </p>
                        {infoVisible && video.videoCount > 1 && (
                            <p className="text-xs text-gray-400 mt-1">{video.videoCount} videos</p>
                        )}
                    </div>
                </div>
            </Link>

            <DropdownMenu>
                <DropdownMenuTrigger className="rounded hover:bg-gray-200 flex-shrink-0 cursor-pointer">
                    <EllipsisVertical className="w-4 h-4 text-gray-500" />
                </DropdownMenuTrigger>
                <DropdownMenuContent align="end" className="min-w-44">
                    <DropdownMenuItem>
                        <a
                            href={`${API_URL}/download/${video.id}`}
                            download={`${video.name}.mp4`}
                            onClick={() => { saveDownloadInfo(video.id, `${video.name}.mp4`); setDownloadInfo(getDownloadInfo(video.id)) }}
                            className="contents"
                        >
                            <Download className="w-4 h-4 shrink-0" />
                            {downloadInfo ? 'Download again' : 'Download'}
                        </a>
                    </DropdownMenuItem>
                    {downloadInfo && (
                        <DropdownMenuItem onClick={() => { localStorage.removeItem(`download_${video.id}`); setDownloadInfo(null) }}>
                            <X className="w-4 h-4 shrink-0" />
                            Clear
                        </DropdownMenuItem>
                    )}
                </DropdownMenuContent>
            </DropdownMenu>
        </div>
    )
}
