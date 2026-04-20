import { forwardRef, useImperativeHandle, useRef, useState } from 'react'

async function copyToClipboard(text, textareaEl) {
    if (navigator.clipboard) {
        await navigator.clipboard.writeText(text)
    } else {
        textareaEl.select()
        document.execCommand('copy')
    }
}
import FormDialog from '@/components/frontend/FormDialog/FormDialog'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'

const WatchedTimesDialogs = forwardRef(function WatchedTimesDialogs({ onImport }, ref) {
    const exportDialogRef = useRef(null)
    const importDialogRef = useRef(null)
    const exportTextareaRef = useRef(null)
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
                onSubmit={async () => { await copyToClipboard(exportData, exportTextareaRef.current) }}
            >
                <Label htmlFor="export-data">Data</Label>
                <Textarea
                    ref={exportTextareaRef}
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
