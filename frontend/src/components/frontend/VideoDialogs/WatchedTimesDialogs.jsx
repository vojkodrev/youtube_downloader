import { forwardRef, useImperativeHandle, useRef, useState } from 'react'
import FormDialog from '@/components/frontend/FormDialog/FormDialog'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'

const WatchedTimesDialogs = forwardRef(function WatchedTimesDialogs({ onImport }, ref) {
    const exportDialogRef = useRef(null)
    const importDialogRef = useRef(null)
    const [exportData, setExportData] = useState('')

    useImperativeHandle(ref, () => ({
        openExport: (data) => { setExportData(data); exportDialogRef.current.open() },
        openImport: () => importDialogRef.current.open(),
    }))

    return (
        <>
            <FormDialog
                ref={exportDialogRef}
                title="Export Watched Times"
                description="Copy the watched times data below."
                submitLabel="Copy"
                onSubmit={async () => { await navigator.clipboard.writeText(exportData) }}
            >
                <Label htmlFor="export-data">Data</Label>
                <Textarea
                    id="export-data"
                    name="data"
                    className="h-64 font-mono text-xs break-all overflow-y-auto resize-none"
                    readOnly
                    value={exportData}
                    onChange={() => {}}
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
                    className="h-64 font-mono text-xs break-all overflow-y-auto resize-none"
                    placeholder="Paste data here..."
                />
            </FormDialog>
        </>
    )
})

export default WatchedTimesDialogs
