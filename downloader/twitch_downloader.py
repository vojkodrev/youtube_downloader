import asyncio
import os

import yt_dlp

from injector import inject

from config import Config
from downloader import Downloader


class TwitchDownloader(Downloader):
    @inject
    def __init__(self, config: Config):
        self._config = config

    async def download(self, url: str) -> None:
        if not url:
            raise ValueError("url is required")

        def sync():
            ydl_opts = {
                "format": "best",
                "merge_output_format": "mp4",
                "overwrites": True,
                "outtmpl": os.path.join(
                    self._config["output_folder"],
                    "[%(uploader)s] %(title)s.%(ext)s",
                ),
                # --- The Retry Trio ---
                "retries": 10,  # Generic network retries
                "fragment_retries": 10,  # Video chunk/fragment retries
                "extractor_retries": 10,  # Website parsing/scraping retries
                "file_access_retries": 10,  # Local disk/NAS access retries
                # --- Precise Timing Control ---
                "sleep_interval": 15,  # Seconds to wait between download tasks
                "max_sleep_interval": 15,  # Keep it strictly at 15s (no randomization)
                "sleep_requests": 5,  # Wait 5s between finding info for each video
                # --- Safety Buffers ---
                "socket_timeout": 30,  # Wait 30s before considering a socket "dead"
            }
            with yt_dlp.YoutubeDL(ydl_opts) as ydl:
                ydl.download([url])

        loop = asyncio.get_running_loop()
        await loop.run_in_executor(None, sync)
