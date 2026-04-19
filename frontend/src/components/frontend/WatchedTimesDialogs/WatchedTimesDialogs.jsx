import { forwardRef, useImperativeHandle, useRef } from 'react'
import FormDialog from '@/components/frontend/FormDialog/FormDialog'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'

const WatchedTimesDialogs = forwardRef(function WatchedTimesDialogs({ onExport, onImport }, ref) {
    const exportDialogRef = useRef(null)
    const importDialogRef = useRef(null)

    useImperativeHandle(ref, () => ({
        openExport: () => exportDialogRef.current.open(),
        openImport: () => importDialogRef.current.open(),
    }))

    return (
        <>
            <FormDialog
                ref={exportDialogRef}
                title="Export Watched Times"
                description="Copy the watched times data below."
                submitLabel="Close"
                onSubmit={onExport}
            >
                <Label htmlFor="export-data">Data</Label>
                <Textarea
                    id="export-data"
                    name="data"
                    className="min-h-64 font-mono text-xs"
                    readOnly
                />
            </FormDialog>

            <FormDialog
                ref={importDialogRef}
                title="Import Watched Times"
                description="Paste watched times data below to import."
                submitLabel="Import"
                onSubmit={onImport}
            >
                <Label htmlFor="import-data">Data</Label>
                <Textarea
                    id="import-data"
                    name="data"
                    className="min-h-64 font-mono text-xs"
                    placeholder="Paste data here..."
                />
            </FormDialog>
        </>
    )
})

export default WatchedTimesDialogs
