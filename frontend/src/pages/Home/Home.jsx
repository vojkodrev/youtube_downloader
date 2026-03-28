export default function Home() {
    return (
        <div className="flex flex-col h-screen">

            {/* Top */}
            <div className="bg-gray-200 p-4">
                <p>Top</p>
            </div>

            {/* Middle: main content + sidebar */}
            <div className="flex flex-1">

                {/* Left: video content + info panel */}
                <div className="flex flex-col flex-1">
                    <div className="h-[50rem] bg-white p-4">
                        <p>Video content</p>
                    </div>
                    <div className="h-[125rem] bg-gray-100 p-4">
                        <p>Info panel</p>
                    </div>
                </div>

                {/* Sidebar */}
                <div className="w-80 bg-gray-50 p-4">
                    <p>Sidebar</p>
                </div>

            </div>
        </div>
    )
}
