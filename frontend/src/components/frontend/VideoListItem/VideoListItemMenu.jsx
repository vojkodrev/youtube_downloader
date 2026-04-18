import { EllipsisVertical, Download, RotateCcw, CheckCheck, Trash2 } from 'lucide-react'
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuSeparator,
    DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { mediaUrl } from '../../../api'

export default function VideoListItemMenu({ video, onWatchedMark, onWatchedReset, onDelete }) {
    return (
        <DropdownMenu>
            <DropdownMenuTrigger className="rounded hover:bg-gray-200 flex-shrink-0 cursor-pointer">
                <EllipsisVertical className="w-4 h-4 text-gray-500" />
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end" className="min-w-50">
                <DropdownMenuItem disabled={video.status !== 'Ready'}>
                    <a
                        href={mediaUrl(`/download/${video.id}`)}
                        download={`${video.name}.mp4`}
                        className="w-full flex items-center gap-2"
                    >
                        <Download className="w-4 h-4 shrink-0" />
                        Download Original
                    </a>
                </DropdownMenuItem>
                {[...(video.versions ?? [])].sort((a, b) => (parseInt(b.quality) || 0) - (parseInt(a.quality) || 0)).map(v => (
                    <DropdownMenuItem key={v.id}>
                        <a href={mediaUrl(`/download/${v.id}`)} download={v.filename} className="w-full flex items-center gap-2">
                            <Download className="w-4 h-4 shrink-0" />
                            Download {v.quality}
                        </a>
                    </DropdownMenuItem>
                ))}
                <DropdownMenuSeparator />
                <DropdownMenuItem
                    onClick={() => onWatchedMark?.(video)}
                >
                    <CheckCheck className="w-4 h-4 shrink-0" />
                    Mark as Watched
                </DropdownMenuItem>
                <DropdownMenuItem
                    onClick={() => onWatchedReset?.(video)}
                >
                    <RotateCcw className="w-4 h-4 shrink-0" />
                    Reset Watched
                </DropdownMenuItem>
                <DropdownMenuSeparator />
                <DropdownMenuItem
                    disabled={video.status !== 'Ready'}
                    className="text-red-600 focus:text-red-600"
                    onClick={() => onDelete?.(video)}
                >
                    <Trash2 className="w-4 h-4 shrink-0" />
                    Delete
                </DropdownMenuItem>
            </DropdownMenuContent>
        </DropdownMenu>
    )
}
