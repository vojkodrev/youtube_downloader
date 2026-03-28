export default function Home() {
    return (
        <div className="flex flex-col h-[125rem]">

            {/* Top */}
            <div className="bg-gray-200 p-4">
                <p>Top</p>
            </div>

            {/* Middle */}
            <div className="flex flex-col md:flex-row flex-1">

                {/* Left: video content + info panel */}
                <div className="flex flex-col md:flex-1">
                    <div className="h-[25rem] md:h-[45rem] bg-white p-4">
                        <p>Video content</p>
                    </div>
                    <div className="h-[10rem] bg-gray-100 p-4">
                        <p>Info panel</p>
                    </div>
                </div>

                {/* Sidebar */}
                <div className="h-[31.25rem] md:h-auto md:w-80 bg-gray-50 p-4">
                    <p>Sidebar</p>
                </div>

            </div>
        </div>
    )
}
