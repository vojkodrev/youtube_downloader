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

const DEFAULT_ERROR = 'An error occurred. Please try again.'

const FormDialog = forwardRef(function FormDialog({ title, description, submitLabel = 'Submit', onSubmit, children }, ref) {
    const [open, setOpen] = useState(false)
    const [error, setError] = useState(null)
    const [pending, setPending] = useState(false)

    useImperativeHandle(ref, () => ({
        open() { setOpen(true) }
    }))

    function handleClose() {
        setError(null)
        setOpen(false)
    }

    async function handleSubmit(e) {
        e.preventDefault()
        setPending(true)
        setError(null)
        try {
            await onSubmit(new FormData(e.target))
            handleClose()
        } catch (err) {
            setError(err.message || DEFAULT_ERROR)
        } finally {
            setPending(false)
        }
    }

    return (
        <Dialog open={open} onOpenChange={open => { if (!open) handleClose() }}>
            <DialogContent>
                <DialogHeader>
                    <DialogTitle>{title}</DialogTitle>
                    {description && <DialogDescription>{description}</DialogDescription>}
                </DialogHeader>
                <form onSubmit={handleSubmit}>
                    <div className="flex flex-col gap-2 py-2">
                        {children}
                        {error && <p className="text-sm text-red-500">{error}</p>}
                    </div>
                    <DialogFooter>
                        <Button type="button" variant="outline" onClick={handleClose}>Cancel</Button>
                        <Button type="submit" disabled={pending}>{submitLabel}</Button>
                    </DialogFooter>
                </form>
            </DialogContent>
        </Dialog>
    )
})

export default FormDialog
