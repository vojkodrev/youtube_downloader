export const VIDEO_SAVE_INTERVAL = 5

export function isFullyWatched(video) {
    const saved = parseFloat(video.savedTime) || 0
    return !!video.duration && saved >= video.duration - VIDEO_SAVE_INTERVAL
}
