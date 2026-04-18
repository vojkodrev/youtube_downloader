import { useState, forwardRef, useImperativeHandle } from 'react'
import {
    Dialog,
    DialogContent,
    DialogHeader,
    DialogTitle,
    DialogDescription,
    DialogFooter,
} from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'

const DeleteVideoDialog = forwardRef(function DeleteVideoDialog(_, ref) {
    const [open, setOpen] = useState(false)
    const [videos, setVideos] = useState([])
    const [onConfirm, setOnConfirm] = useState(null)

    useImperativeHandle(ref, () => ({
        open(videos, confirmDelete) { setVideos(videos); setOnConfirm(() => confirmDelete); setOpen(true) }
    }))

    function handleClose() {
        setOpen(false)
    }

    function handleConfirm() {
        setOpen(false)
        onConfirm(videos)
    }

    return (
        <Dialog open={open} onOpenChange={open => { if (!open) handleClose() }}>
            <DialogContent>
                <DialogHeader>
                    <DialogTitle>Delete Video</DialogTitle>
                    <DialogDescription>
                        {videos.length > 1
                            ? `Are you sure you want to delete all ${videos.length} videos in this playlist? This cannot be undone.`
                            : 'Are you sure you want to delete this video? This cannot be undone.'
                        }
                    </DialogDescription>
                </DialogHeader>
                <DialogFooter>
                    <Button variant="outline" onClick={handleClose}>Cancel</Button>
                    <Button variant="destructive" onClick={handleConfirm}>Delete</Button>
                </DialogFooter>
            </DialogContent>
        </Dialog>
    )
})

export default DeleteVideoDialog
